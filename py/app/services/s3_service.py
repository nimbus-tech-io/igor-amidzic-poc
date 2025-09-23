import logging
import boto3
from botocore.exceptions import ClientError
from app.core.config import config

logger = logging.getLogger(__name__)


class S3Service:
    def __init__(self):
        self.client = None
        self._initialize_client()
    
    def _initialize_client(self):
        """Initialize S3 client."""
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
        logger.info("S3 client initialized")
    
    def read_file(self, file_name: str) -> str:
        """Read file content from S3 bucket."""
        try:
            logger.info(f"Reading file {file_name} from S3 bucket {config.S3_BUCKET_NAME}")
            
            response = self.client.get_object(
                Bucket=config.S3_BUCKET_NAME,
                Key=file_name
            )
            
            file_content = response['Body'].read().decode('utf-8')
            logger.info(f"Successfully read file {file_name}, size: {len(file_content)} chars")
            return file_content
            
        except ClientError as e:
            logger.error(f"Failed to read file {file_name} from S3: {e}")
            raise
        except Exception as e:
            logger.error(f"Unexpected error reading file {file_name}: {e}")
            raise
    
    def file_exists(self, file_name: str) -> bool:
        """Check if file exists in S3 bucket."""
        try:
            self.client.head_object(Bucket=config.S3_BUCKET_NAME, Key=file_name)
            return True
        except ClientError:
            return False


s3_service = S3Service()