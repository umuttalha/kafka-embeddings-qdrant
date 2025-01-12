import os
import json
from dotenv import load_dotenv
from kafka import KafkaConsumer, KafkaProducer
from sentence_transformers import SentenceTransformer
from qdrant_client import QdrantClient
from qdrant_client.models import Distance, VectorParams, PointStruct

load_dotenv()

def init_kafka_producer():
    kafka_host = os.getenv('KAFKA_HOST', 'localhost')
    return KafkaProducer(
        bootstrap_servers=[f'{kafka_host}:9092'],
        value_serializer=lambda x: json.dumps(x).encode('utf-8')
    )

def search_similar(query, model, client, limit=3):
    try:
        # Create embedding for query
        query_vector = model.encode(query)
        
        # Search in Qdrant
        search_results = client.search(
            collection_name='text_embeddings',
            query_vector=query_vector,
            limit=limit
        )
        
        # Extract results
        results = [
            {
                "text": hit.payload["text"],
                "score": float(hit.score)  # Convert numpy float to Python float
            }
            for hit in search_results
        ]
        
        return {"status": "success", "results": results}
    except Exception as e:
        print(f"Error searching: {e}")
        return {"status": "error", "message": str(e)}

def main():
    kafka_host = os.getenv('KAFKA_HOST', 'localhost')
    qdrant_host = os.getenv('QDRANT_HOST', 'localhost')
    
    # Initialize Kafka consumer and producer
    consumer = KafkaConsumer(
        'text-topic',
        bootstrap_servers=[f'{kafka_host}:9092'],
        value_deserializer=lambda x: json.loads(x.decode('utf-8'))
    )
    producer = init_kafka_producer()
    
    # Initialize model and client
    model = SentenceTransformer('all-MiniLM-L6-v2')
    client = QdrantClient(host=qdrant_host, port=6333)
    
    for message in consumer:
        data = message.value
        
        if data['type'] == 'insert':
            try:
                embedding = model.encode(data['content'])
                client.upsert(
                    collection_name='text_embeddings',
                    points=[PointStruct(
                        vector=embedding.tolist(),
                        payload={'text': data['content']}
                    )]
                )
            except Exception as e:
                print(f"Error inserting: {e}")
                
        elif data['type'] == 'search':
            try:
                results = search_similar(data['content'], model, client)
                # Send results back through Kafka
                producer.send('response-topic', {
                    'request_id': data.get('request_id'),
                    'results': results
                })
                print(f"Search results sent: {results}")
            except Exception as e:
                print(f"Error searching: {e}")

if __name__ == "__main__":
    main()