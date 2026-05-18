# Recommendation Engine clickstream & Processing Pipeline

The personalized recommendation module increases session duration and Conversion Rates.

```mermaid
graph LR
    Client[Next.js Client] -->|Clickstream Events| Kafka[Kafka Event Broker]
    Kafka -->|Stream ingestion| Flink[Apache Flink Analytics]
    Flink -->|Sliding Window Aggregation| Features[(Redis Feature Store)]
    Features -->|Feature Vector| RankModel[ML Ranking Model: Python]
    RankModel -->|Recommendation list| Client
```

## Apache Flink clickstream Aggregation
- Flink reads click/view counts across a 5-minute sliding window to identify trending SKU groups.
- Product listings retrieve these scores to adjust the ranking of user homepages dynamically.
