import uvicorn
import logging
from app.core.app import create_app
from app.core.config import config


uvicorn_logger = logging.getLogger("uvicorn")
uvicorn_logger.setLevel(logging.INFO)

app = create_app()

if __name__ == "__main__":
    print("🎯 Starting FastAPI application...")
    print("📋 Logging configured - you should see logs in console!")
    uvicorn.run(
        "main:app",
        host=config.HOST,
        port=config.PORT,
        reload=config.DEBUG,
        log_level="info"
    )