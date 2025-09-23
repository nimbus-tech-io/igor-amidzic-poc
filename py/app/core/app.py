import asyncio
from contextlib import asynccontextmanager
from fastapi import FastAPI
from app.api.routes import router
from app.services.sqs_worker import sqs_worker
from app.core.config import config


@asynccontextmanager
async def lifespan(app: FastAPI):
    """Lifespan manager for FastAPI app - handles startup and shutdown."""
    print("Starting Py application...")
    
    worker_task = asyncio.create_task(sqs_worker.start())
    
    yield
    
    print("🛑 Shutting down Py application...")
    sqs_worker.stop()
    worker_task.cancel()
    
    try:
        await worker_task
    except asyncio.CancelledError:
        print("✅ SQS Worker stopped successfully")


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