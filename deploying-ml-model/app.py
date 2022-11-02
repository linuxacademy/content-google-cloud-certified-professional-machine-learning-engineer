import os
import json
 
from flask import Flask, request, Response
import torch
from transformers import (
   DistilBertTokenizer, DistilBertForQuestionAnswering,
)
 
class Model(object):
   def __init__(self, tokenizer, model):
       self.tokenizer = tokenizer
       self.model = model
 
   def encode(self, question, context):
       encoded = self.tokenizer.encode_plus(question, context)
       return encoded["input_ids"], encoded["attention_mask"]
 
   def decode(self, token):
       answer_tokens = self.tokenizer.convert_ids_to_tokens(token , skip_special_tokens=True)
       return self.tokenizer.convert_tokens_to_string(answer_tokens)
 
   def predict(self, input):
       question, context = input['question'], input['context']
       input_ids, attention_mask = self.encode(question, context)
       start_scores, end_scores = self.model(torch.tensor([input_ids]), attention_mask=torch.tensor([attention_mask]))
       ans_tokens = input_ids[torch.argmax(start_scores) : torch.argmax(end_scores)+1]
       answer = self.decode(ans_tokens)
       return answer
 
# DistilBERT
model = Model(
   tokenizer=DistilBertTokenizer.from_pretrained('./model', return_token_type_ids=True),
   model=DistilBertForQuestionAnswering.from_pretrained('./model'),
)
 
app = Flask(__name__)
 
@app.route('/')
def hello_world():
   target = os.environ.get('TARGET', 'World')
   return 'Hello {}!\n'.format(target)
 
@app.route('/predict', methods=['POST'])
def predict():
   sentiment = request.get_json()
   return Response(json.dumps(model.predict(sentiment)),  mimetype='application/json')
  
if __name__ == "__main__":
   app.run(debug=True,host='0.0.0.0',port=int(os.environ.get('PORT', 8080)))
