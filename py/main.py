import uvicorn
from app.core.app import create_app
from app.core.config import config

app = create_app()

if __name__ == "__main__":
    uvicorn.run(
        "main:app",
        host=config.HOST,
        port=config.PORT,
        reload=config.DEBUG
    )