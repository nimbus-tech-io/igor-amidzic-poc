import asyncio
import logging
import boto3
from botocore.exceptions import ClientError
from app.core.config import config

logger = logging.getLogger(__name__)


class SQSWorker:
    def __init__(self):
        self.running = False
        self.sqs_client = None
    
    async def start(self):
        """Start the SQS worker."""
        logger.info("Starting SQS Worker...")
        
        client_config = {
            'aws_access_key_id': config.AWS_ACCESS_KEY_ID,
            'aws_secret_access_key': config.AWS_SECRET_ACCESS_KEY,
            'region_name': config.AWS_REGION
        }
        if config.AWS_ENDPOINT:
            client_config['endpoint_url'] = config.AWS_ENDPOINT
        
        self.sqs_client = boto3.client('sqs', **client_config)
        self.running = True
        
        # Poll messages
        while self.running:
            try:
                response = self.sqs_client.receive_message(
                    QueueUrl=config.SQS_QUEUE_URL,
                    MaxNumberOfMessages=10,
                    WaitTimeSeconds=20
                )
                
                messages = response.get('Messages', [])
                if messages:
                    logger.info(f"Received {len(messages)} message(s)")
                    for message in messages:
                        await self.process_message(message)
                
            except ClientError as e:
                logger.error(f"SQS error: {e}")
                await asyncio.sleep(5)
            except Exception as e:
                logger.error(f"Unexpected error: {e}")
                await asyncio.sleep(5)
    
    async def process_message(self, message):
        """Process a single message."""
        logger.info(f"Processing: {message.get('Body', '')}")
        
        await asyncio.sleep(1)  
        
        self.sqs_client.delete_message(
            QueueUrl=config.SQS_QUEUE_URL,
            ReceiptHandle=message['ReceiptHandle']
        )
        logger.info("Message processed and deleted")
    
    def stop(self):
        """Stop the worker."""
        self.running = False


sqs_worker = SQSWorker()

