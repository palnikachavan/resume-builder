package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// Resume represents a user's resume
type Resume struct {
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Skills     []string  `json:"skills"`
	Experience []string  `json:"experience"`
	Projects   []Project `json:"projects"`
}

// Project represents a project
type Project struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	URL         string  `json:"url"`
	Score       float64 `json:"score"`
}

// ProjectRecommendationInput represents input for project recommendations
type ProjectRecommendationInput struct {
	Role string `json:"role"`
	TopN int    `json:"top_n"`
}

// In-memory storage (use a DB in production)
var resumes = make(map[string]Resume)
var projects []Project

// Load projects from a JSON file
func loadProjects() {
	file, err := os.Open("projects.json")
	if err != nil {
		log.Println("No project data found, starting fresh.")
		return
	}
	defer file.Close()

	err = json.NewDecoder(file).Decode(&projects)
	if err != nil {
		log.Println("Error decoding projects JSON:", err)
	}
}

// Save resumes to a JSON file
func saveResumes() {
	file, err := os.Create("resumes.json")
	if err != nil {
		log.Println("Error saving resumes:", err)
		return
	}
	defer file.Close()

	json.NewEncoder(file).Encode(resumes)
}

// CreateResume handles creating a new resume
func CreateResume(c *gin.Context) {
	var resume Resume
	if err := c.ShouldBindJSON(&resume); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if _, exists := resumes[resume.Email]; exists {
		c.JSON(400, gin.H{"error": "Resume already exists"})
		return
	}

	resumes[resume.Email] = resume
	saveResumes()

	c.JSON(201, gin.H{"message": "Resume created successfully"})
}

// GetResume handles retrieving a resume by email
func GetResume(c *gin.Context) {
	email := c.Param("email")
	resume, exists := resumes[email]
	if !exists {
		c.JSON(404, gin.H{"error": "Resume not found"})
		return
	}

	c.JSON(200, resume)
}

// RecommendProjects suggests projects based on a role
func RecommendProjects(c *gin.Context) {
	var input ProjectRecommendationInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	if len(projects) == 0 {
		c.JSON(500, gin.H{"error": "No projects available"})
		return
	}

	// Simple scoring: Check keyword match
	roleWords := strings.Fields(strings.ToLower(input.Role))
	var recommended []Project
	for _, proj := range projects {
		for _, word := range roleWords {
			if strings.Contains(strings.ToLower(proj.Description), word) {
				proj.Score = rand.Float64() * 10 // Assign a random score for now
				recommended = append(recommended, proj)
				break
			}
		}
	}

	// Limit results
	if len(recommended) > input.TopN {
		recommended = recommended[:input.TopN]
	}

	c.JSON(200, gin.H{"recommended_projects": recommended})
}

// AddProjectToResume adds a project to a resume
func AddProjectToResume(c *gin.Context) {
	email := c.Param("email")
	var project Project

	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request body"})
		return
	}

	resume, exists := resumes[email]
	if !exists {
		c.JSON(404, gin.H{"error": "Resume not found"})
		return
	}

	resume.Projects = append(resume.Projects, project)
	resumes[email] = resume
	saveResumes()

	c.JSON(200, gin.H{"message": "Project added to resume"})
}

func main() {
	loadProjects()

	r := gin.Default()
	r.POST("/create-resume", CreateResume)
	r.GET("/get-resume/:email", GetResume)
	r.POST("/recommend-projects", RecommendProjects)
	r.POST("/add-project-to-resume/:email", AddProjectToResume)

	fmt.Println("Server running on http://localhost:8080")
	r.Run(":8080") // Start server on port 8080
}
