import asyncio
from contextlib import asynccontextmanager
from fastapi import FastAPI
from app.api.routes import router
from app.worker.worker import worker
from app.core.config import config


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Lifespan manager for FastAPI app - handles startup and shutdown."""
    
    worker_task = asyncio.create_task(worker.start())
    print("✅ Worker task created")
        
    yield
    
    print("🛑 Shutting down Py application...")
    worker.stop()
    worker_task.cancel()
    
    try:
        await worker_task
    except asyncio.CancelledError:
        print("✅ Worker stopped successfully")


def create_app() -> FastAPI:
    """Create and configure FastAPI application."""
    app = FastAPI(
        title="FastAPI SQS Worker App",
        description="Basic FastAPI application with AWS SQS worker",
        version="1.0.0",
        debug=config.DEBUG,
        lifespan=lifespan
    )
    
    app.include_router(router)
    
    return app