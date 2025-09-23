from fastapi import APIRouter
from app.services.worker import worker
from app.core.config import config

router = APIRouter()


@router.get("/")
async def root():
    """Root endpoint."""
    return {
        "message": "FastAPI SQS Worker App is running!",
        "status": "healthy"
    }


@router.get("/health")
async def health_check():
    """Health check endpoint."""
    return {
        "status": "healthy",
        "worker_running": worker.running,
        "queue_url": config.SQS_QUEUE_URL,
        "environment": config._config_data.get("environment"),
        "aws_endpoint": config.AWS_ENDPOINT
    }


@router.get("/worker/status")
async def worker_status():
    """Get worker status."""
    return {
        "worker_running": worker.running,
        "queue_url": config.SQS_QUEUE_URL,
        "queue_name": config.SQS_QUEUE_NAME,
        "s3_bucket": config.S3_BUCKET_NAME,
        "aws_configured": config.validate_aws_config(),
        "aws_region": config.AWS_REGION,
        "aws_endpoint": config.AWS_ENDPOINT,
        "environment": config._config_data.get("environment")
    }