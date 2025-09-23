import logging
import boto3
import base64
from botocore.exceptions import ClientError
from app.core.config import config

logger = logging.getLogger(__name__)


class S3Service:
    def __init__(self):
        self.client = None
        self._initialize_client()
    
    def _initialize_client(self):
        """Initialize S3 client."""
        try:
            client_config = {
                'aws_access_key_id': config.AWS_ACCESS_KEY_ID,
                'aws_secret_access_key': config.AWS_SECRET_ACCESS_KEY,
                'region_name': config.AWS_REGION
            }
            if config.AWS_ENDPOINT:
                client_config['endpoint_url'] = config.AWS_ENDPOINT
            
            if config.S3_FORCE_PATH_STYLE:
                client_config['config'] = boto3.session.Config(s3={'addressing_style': 'path'})
            
            self.client = boto3.client('s3', **client_config)
        except Exception as e:
            print(f"❌ S3 client init failed: {e}")
            self.client = None
    
    def read_file(self, file_name: str) -> str:
        """Read file from S3 and return as base64."""
        if not self.client:
            raise Exception("S3 client not configured")
            
        try:
            response = self.client.get_object(
                Bucket=config.S3_BUCKET_NAME,
                Key=file_name
            )
            
            raw_data = response['Body'].read()
            base64_content = base64.b64encode(raw_data).decode('utf-8')
            
            return base64_content
            
        except ClientError as e:
            raise Exception(f"Failed to read file {file_name}: {e}")
        except Exception as e:
            raise Exception(f"Error reading file {file_name}: {e}")
    
    def file_exists(self, file_name: str) -> bool:
        """Check if file exists in S3."""
        if not self.client:
            return False
            
        try:
            self.client.head_object(Bucket=config.S3_BUCKET_NAME, Key=file_name)
            return True
        except ClientError:
            return False


s3_service = S3Service()