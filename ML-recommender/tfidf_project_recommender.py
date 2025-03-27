from fastapi import FastAPI, HTTPException
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.metrics.pairwise import cosine_similarity
import json
import uvicorn
import os

app = FastAPI()

json_file_path = "repositories.json"

# Load repositories safely
repositories = []
if os.path.exists(json_file_path):
    try:
        with open(json_file_path, "r") as file:
            data = file.read()
            if data.strip():  # Ensure file isn't empty
                repositories = json.loads(data)
            else:
                print("Warning: JSON file is empty.")
    except (json.JSONDecodeError, ValueError) as e:
        print(f"Error loading JSON: {e}")

# Ensure we have valid data
if not repositories:
    repositories = [{"name": "Dummy Project", "description": "Sample description", "html_url": "#"}]

# Extract descriptions
repo_descriptions = [repo.get("description", "") for repo in repositories]
repo_names = [repo.get("name", "Unknown") for repo in repositories]
repo_urls = [repo.get("html_url", "#") for repo in repositories]

# Initialize TF-IDF Vectorizer
vectorizer = TfidfVectorizer(stop_words="english")
repo_vectors = vectorizer.fit_transform(repo_descriptions) if repo_descriptions else None

@app.post("/recommend-projects")
def recommend_projects(role: str, top_n: int = 5):
    if not role:
        raise HTTPException(status_code=400, detail="Role description required")
    
    if repo_vectors is None:
        raise HTTPException(status_code=500, detail="No valid repository data available")
    
    role_vector = vectorizer.transform([role])
    similarities = cosine_similarity(role_vector, repo_vectors).flatten()
    ranked_indices = similarities.argsort()[::-1][:top_n]
    
    recommendations = [{
        "name": repo_names[i],
        "description": repo_descriptions[i],
        "url": repo_urls[i],
        "score": float(similarities[i])
    } for i in ranked_indices]
    
    return {"recommended_projects": recommendations}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000, reload=True)
