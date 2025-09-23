import logging
import httpx
from typing import Optional
from app.core.config import config

logger = logging.getLogger(__name__)


class OpenRouterService:
    def __init__(self):
        self.api_key = config.OPENROUTER_API_KEY
        self.model = config.OPENROUTER_MODEL
        self.base_url = config.OPENROUTER_BASE_URL
        
        if self.api_key:
            masked_key = f"{self.api_key[:10]}...{self.api_key[-4:]}" if len(self.api_key) > 14 else "****"
            logger.info(f"🔑 OpenRouter API key loaded: {masked_key}")
            logger.info(f"🤖 Using model: {self.model}")
        else:
            logger.warning("❌ OpenRouter API key not configured!")
            logger.warning("💡 Set OPENROUTER_API_KEY environment variable")
    
    async def analyze_file(self, base64_content: str, file_name: str) -> Optional[str]:
        """
        Send base64 file content to OpenRouter AI for analysis.
        
        Args:
            base64_content: Base64 encoded file content 
            file_name: Name of the file for context
            
        Returns:
            AI-generated overview of the file or None if error
        """
        if not self.api_key:
            logger.error("OpenRouter API key not configured")
            return None
        
        try:
            headers = {
                "Authorization": f"Bearer {self.api_key}",
                "Content-Type": "application/json",
                "HTTP-Referer": "http://localhost:8081",
                "X-Title": "Nimbus File Analyzer"
            }
            
            messages = [
                {
                    "role": "user", 
                    "content": f"""Analyze this file: {file_name}

Please provide:
1. File type and format
2. Key content overview  
3. Main insights or information

Base64 content: {base64_content}
"""
                }
            ]
            
            payload = {
                "model": self.model,
                "messages": messages,
                "max_tokens": 1000,
                "temperature": 0.7
            }

            async with httpx.AsyncClient(timeout=30.0) as client:
                response = await client.post(
                    f"{self.base_url}/chat/completions",
                    headers=headers,
                    json=payload
                )
                
                if response.status_code == 200:
                    result = response.json()
                    overview = result["choices"][0]["message"]["content"]
                    logger.info(f"Successfully analyzed file: {file_name}")
                    return overview
                else:
                    logger.error(f"OpenRouter API error: {response.status_code} - {response.text}")
                    return None
                    
        except httpx.TimeoutException:
            logger.error("OpenRouter API request timed out")
            return None
        except Exception as e:
            logger.error(f"Error calling OpenRouter API: {e}")
            return None

openrouter_service = OpenRouterService()