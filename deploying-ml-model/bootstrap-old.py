import torch
 
from transformers import (
   DistilBertTokenizer, DistilBertForQuestionAnswering,
)
 
DistilBertTokenizer.from_pretrained('distilbert-base-uncased', return_token_type_ids=True).save_pretrained('./model')
DistilBertForQuestionAnswering.from_pretrained('distilbert-base-uncased-distilled-squad').save_pretrained('./model')
