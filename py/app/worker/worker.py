import asyncio
import logging
import json
from botocore.exceptions import ClientError
from app.services.sqs_service import sqs_service
from app.services.s3_service import s3_service

logger = logging.getLogger(__name__)


class Worker:
    def __init__(self):
        self.running = False
    
    async def start(self):
        """Start the worker."""
        logger.info("Starting Worker...")
        self.running = True
        
        while self.running:
            try:
                messages = sqs_service.receive_messages()
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
        try:
            message_body = message.get('Body', '')
            logger.info(f"Processing message: {message_body}")
            
            try:
                message_data = json.loads(message_body)
                file_name = message_data.get('file_name')
                
                if not file_name:
                    logger.error("Message does not contain 'file_name' field")
                    return
                    
            except json.JSONDecodeError:
                file_name = message_body.strip()
                logger.info(f"Treating message body as file_name: {file_name}")
            
            file_content = s3_service.read_file(file_name)
            
            logger.info(f"File content preview: {file_content[:100]}...")
            await asyncio.sleep(1) 
            
            sqs_service.delete_message(message['ReceiptHandle'])
            logger.info(f"Message processed and deleted for file: {file_name}")
            
        except Exception as e:
            logger.error(f"Error processing message: {e}")
    
    def stop(self):
        """Stop the worker."""
        logger.info("Stopping Worker...")
        self.running = False


worker = Worker()