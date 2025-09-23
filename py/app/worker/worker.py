import asyncio
import json
from botocore.exceptions import ClientError
from app.services.sqs_service import sqs_service
from app.services.s3_service import s3_service
from app.services.openrouter_service import openrouter_service


class Worker:
    def __init__(self):
        self.running = False
    
    async def start(self):
        """Start the worker."""
        print("� Worker started")
        self.running = True
        
        while self.running:
            try:
                messages = sqs_service.receive_messages()
                if messages:
                    for message in messages:
                        print("Processing message...") 
                        await self.process_message(message)
                
                await asyncio.sleep(1)
                
            except ClientError as e:
                print(f"❌ SQS error: {e}")
                await asyncio.sleep(5)
            except Exception as e:
                print(f"❌ Worker error: {e}")
                await asyncio.sleep(5)
    
    async def process_message(self, message):
        """Process a single message."""
        try:
            message_body = message.get('Body', '')
            message_data = json.loads(message_body)
            s3_key = message_data.get('s3_key')  
            file_name = message_data.get('file_name')  
            
            if not s3_key:
                print(f"❌ No s3_key in message")
                return
            
            file_content = s3_service.read_file(s3_key)
            if not file_content:
                print(f"❌ Empty file: {s3_key}")
                return
            
            overview = await openrouter_service.analyze_file(file_content, file_name)
            if overview:
                print(f"\n🤖 AI Analysis for {file_name}:")
                print(f"{overview}\n")
            else:
                print(f"❌ AI analysis failed for: {file_name}")
            
            sqs_service.delete_message(message['ReceiptHandle'])
            
        except Exception as e:
            print(f"❌ Error processing message: {e}")
    
    def stop(self):
        """Stop the worker."""
        print("🛑 Worker stopped")
        self.running = False


worker = Worker()