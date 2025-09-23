import json
import os
from pathlib import Path
from typing import Optional, Dict, Any


class Config:
    """Configuration class for the FastAPI application and SQS worker."""
    
    def __init__(self):
        self._config_data = self._load_config()
        
    def _load_config(self) -> Dict[str, Any]:
        """Load configuration from config.json file in project root."""
        config_path = Path(__file__).parent.parent.parent.parent / "config.json"
        
        try:
            with open(config_path, 'r') as f:
                config_data = json.load(f)
            return config_data
        except FileNotFoundError:
            raise FileNotFoundError(f"Config file not found at {config_path}")
        except json.JSONDecodeError as e:
            raise ValueError(f"Invalid JSON in config file: {e}")
    
    def _resolve_env_var(self, value: str) -> str:
        """Resolve environment variable placeholders like ${VAR_NAME:default}."""
        if isinstance(value, str) and value.startswith("${") and value.endswith("}"):
            var_content = value[2:-1]  
            if ":" in var_content:
                var_name, default_value = var_content.split(":", 1)
                return os.getenv(var_name, default_value)
            else:
                var_name = var_content
                return os.getenv(var_name, "")
        return value
    
    @property
    def AWS_ACCESS_KEY_ID(self) -> Optional[str]:
        return self._config_data.get("aws", {}).get("credentials", {}).get("access_key")
    
    @property
    def AWS_SECRET_ACCESS_KEY(self) -> Optional[str]:
        return self._config_data.get("aws", {}).get("credentials", {}).get("secret_key")
    
    @property
    def AWS_REGION(self) -> str:
        return self._config_data.get("aws", {}).get("region", "us-east-1")
    
    @property
    def AWS_ENDPOINT(self) -> Optional[str]:
        """AWS endpoint for LocalStack or custom endpoint."""
        return self._config_data.get("aws", {}).get("endpoint")
    
    @property
    def SQS_QUEUE_NAME(self) -> Optional[str]:
        return self._config_data.get("sqs", {}).get("queue_name")
    
    @property
    def SQS_QUEUE_URL(self) -> Optional[str]:
        """Build SQS queue URL from configuration."""
        queue_name = self.SQS_QUEUE_NAME
        if not queue_name:
            return None            
        return f"{self.AWS_ENDPOINT}/000000000000/{queue_name}"
            
    @property
    def HOST(self) -> str:
        return os.getenv("HOST", "0.0.0.0")
    
    @property
    def PORT(self) -> int:
        config_port = self._config_data.get("servers", {}).get("python_api", {}).get("port", ":8081")
        if config_port.startswith(":"):
            config_port = config_port[1:] 
        
        return int(os.getenv("PORT", config_port))
    
    @property
    def DEBUG(self) -> bool:
        env_debug = os.getenv("DEBUG", "").lower()
        if env_debug:
            return env_debug == "true"
        
        return self._config_data.get("environment") == "development"
    
    def validate_aws_config(self) -> bool:
        """Validate that all required AWS configuration is present."""
        required_aws_vars = [
            self.AWS_ACCESS_KEY_ID,
            self.AWS_SECRET_ACCESS_KEY,
            self.SQS_QUEUE_NAME
        ]
        return all(var is not None for var in required_aws_vars)


config = Config()