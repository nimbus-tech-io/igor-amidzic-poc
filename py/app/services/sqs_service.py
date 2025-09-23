import logging
import boto3
from botocore.exceptions import ClientError
from typing import List, Dict, Any
from app.core.config import config

logger = logging.getLogger(__name__)


class SQSService:
    def __init__(self):
        self.client = None
        self._initialize_client()
    
    def _initialize_client(self):
        """Initialize SQS client."""
        client_config = {
            'aws_access_key_id': config.AWS_ACCESS_KEY_ID,
            'aws_secret_access_key': config.AWS_SECRET_ACCESS_KEY,
            'region_name': config.AWS_REGION
        }
        if config.AWS_ENDPOINT:
            client_config['endpoint_url'] = config.AWS_ENDPOINT
        
        self.client = boto3.client('sqs', **client_config)
        logger.info("SQS client initialized")
    
    def receive_messages(self, max_messages: int = 10, wait_time: int = 20) -> List[Dict[str, Any]]:
        """Receive messages from SQS queue."""
        try:
            response = self.client.receive_message(
                QueueUrl=config.SQS_QUEUE_URL,
                MaxNumberOfMessages=max_messages,
                WaitTimeSeconds=wait_time
            )
            return response.get('Messages', [])
        except ClientError as e:
            logger.error(f"Error receiving messages from SQS: {e}")
            raise
    
    def delete_message(self, receipt_handle: str):
        """Delete message from SQS queue."""
        try:
            self.client.delete_message(
                QueueUrl=config.SQS_QUEUE_URL,
                ReceiptHandle=receipt_handle
            )
            logger.info("Message deleted from queue")
        except ClientError as e:
            logger.error(f"Failed to delete message: {e}")
            raise


sqs_service = SQSService()