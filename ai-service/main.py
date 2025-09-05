"""
SeedFlow AI Service - Main FastAPI Application
AIサービスのメインアプリケーション
"""

from fastapi import FastAPI, HTTPException, Depends
from fastapi.middleware.cors import CORSMiddleware
from pydantic import BaseModel
import logging
import os
from typing import Optional, List
import structlog

# Configure structured logging
structlog.configure(
    processors=[
        structlog.stdlib.filter_by_level,
        structlog.stdlib.add_logger_name,
        structlog.stdlib.add_log_level,
        structlog.stdlib.PositionalArgumentsFormatter(),
        structlog.processors.TimeStamper(fmt="iso"),
        structlog.dev.ConsoleRenderer()
    ],
    context_class=dict,
    logger_factory=structlog.stdlib.LoggerFactory(),
    wrapper_class=structlog.stdlib.BoundLogger,
    cache_logger_on_first_use=True,
)

logger = structlog.get_logger()

# Initialize FastAPI app
app = FastAPI(
    title="SeedFlow AI Service",
    description="AI processing service for knowledge management",
    version="1.0.0"
)

# Add CORS middleware
app.add_middleware(
    CORSMiddleware,
    allow_origins=["http://localhost:8080"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)

# Pydantic models for request/response
class HealthResponse(BaseModel):
    status: str
    version: str
    ai_models: List[str]

class ProcessTextRequest(BaseModel):
    text: str
    source_type: str = "text"
    source_url: Optional[str] = None

class ProcessTextResponse(BaseModel):
    title: str
    summary: str
    problem: Optional[str] = None
    solution: Optional[str] = None
    result: Optional[str] = None
    insight: Optional[str] = None
    keywords: List[str]
    category: str
    quality_score: int

@app.get("/")
async def root():
    """Root endpoint"""
    return {"message": "SeedFlow AI Service", "version": "1.0.0"}

@app.get("/ai/health", response_model=HealthResponse)
async def health_check():
    """Health check endpoint"""
    try:
        # Check if required environment variables are set
        openai_key = os.getenv("OPENAI_API_KEY")
        claude_key = os.getenv("CLAUDE_API_KEY")
        
        available_models = []
        if openai_key:
            available_models.append("openai-gpt4")
        if claude_key:
            available_models.append("claude-3-sonnet")
        
        if not available_models:
            raise HTTPException(
                status_code=503, 
                detail="No AI API keys configured"
            )
        
        logger.info("Health check passed", models=available_models)
        
        return HealthResponse(
            status="healthy",
            version="1.0.0",
            ai_models=available_models
        )
    except Exception as e:
        logger.error("Health check failed", error=str(e))
        raise HTTPException(status_code=503, detail=str(e))

@app.post("/ai/process", response_model=ProcessTextResponse)
async def process_text(request: ProcessTextRequest):
    """Process text and generate knowledge structure"""
    try:
        logger.info(
            "Processing text request",
            source_type=request.source_type,
            text_length=len(request.text)
        )
        
        # TODO: Implement actual AI processing
        # For now, return a mock response
        
        # Simple keyword extraction (placeholder)
        keywords = extract_keywords(request.text)
        
        # Mock response based on text content
        response = ProcessTextResponse(
            title=f"Knowledge from {request.source_type}",
            summary=request.text[:200] + "..." if len(request.text) > 200 else request.text,
            problem="Generated problem description",
            solution="Generated solution description",
            result="Generated result description",
            insight="Generated insight",
            keywords=keywords,
            category="その他",
            quality_score=75
        )
        
        logger.info("Text processing completed", title=response.title)
        return response
        
    except Exception as e:
        logger.error("Text processing failed", error=str(e))
        raise HTTPException(status_code=500, detail=f"Processing failed: {str(e)}")

@app.post("/ai/extract-url")
async def extract_from_url(url: str):
    """Extract content from URL"""
    try:
        logger.info("Extracting content from URL", url=url)
        
        # TODO: Implement URL extraction
        # For now, return mock response
        return {
            "title": "Extracted Title",
            "content": "Extracted content from URL",
            "metadata": {
                "url": url,
                "extracted_at": "2024-01-01T00:00:00Z"
            }
        }
        
    except Exception as e:
        logger.error("URL extraction failed", url=url, error=str(e))
        raise HTTPException(status_code=500, detail=f"URL extraction failed: {str(e)}")

def extract_keywords(text: str, max_keywords: int = 10) -> List[str]:
    """Simple keyword extraction (placeholder implementation)"""
    # This is a simple implementation - in practice, use more sophisticated NLP
    words = text.lower().split()
    # Filter out common words and short words
    common_words = {"の", "に", "を", "は", "が", "で", "と", "て", "に", "から", "まで", 
                   "the", "a", "an", "and", "or", "but", "in", "on", "at", "to", "for"}
    
    keywords = [word.strip(".,!?;:") for word in words 
               if len(word) > 2 and word.lower() not in common_words]
    
    # Return unique keywords, limited to max_keywords
    return list(set(keywords))[:max_keywords]

if __name__ == "__main__":
    import uvicorn
    
    # Configure logging
    logging.basicConfig(level=logging.INFO)
    
    # Start the server
    uvicorn.run(
        "main:app",
        host="0.0.0.0",
        port=8001,
        reload=False,
        log_level="info"
    )