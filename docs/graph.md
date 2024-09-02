# Graphs in OpenAgents

The mobile app's main screen displays a graph representing the agent's knowledge base, which is fundamental to OpenAgents' information processing and retrieval.

## Graph Structure and Data Representation

The knowledge base uses graph theory concepts:
- Nodes represent discrete pieces of information or concepts
- Edges represent relationships between nodes
- Nodes can be grouped into communities or topics, each with its own summary

Each node has a vector embedding, enabling cosine similarity search and specific information retrieval. This approach, inspired by Microsoft's GraphRAG paper, combines graph-based representation with vector-based similarity search.

## GraphRAG Implementation

OpenAgents implements GraphRAG (Graph-based Retrieval Augmented Generation) to improve understanding of complex relationships within codebases and user data. Key aspects include:

1. **Code Understanding**: Creating a knowledge graph of the codebase
2. **Query Responses**: Providing contextually relevant responses to user queries
3. **User Data Analysis**: Comprehensively understanding user interactions over time
4. **Advanced Search**: Enabling sophisticated search functionality

The implementation process involves defining text units, extracting entities and relationships, constructing a knowledge graph, generating summaries, and enhancing the query system.

## Nostr-based Communication

OpenAgents uses Nostr for communication between the mobile app (client) and the relay (distributed backend), leveraging custom event kinds:

1. Agent Command Request (kind 5838)
2. Agent Command Response (kind 6838)
3. Speech-to-Text Request (kind 5252)
4. Speech-to-Text Response (kind 6252)

### Communication Flow

1. User initiates an action in the app
2. App sends relevant events to the relay (e.g., speech-to-text, command request)
3. Relay processes the request, potentially involving graph-building (indexing)
4. Relay sends update events and final response back to the app

This event-based system enables asynchronous, real-time communication between the app and the backend.

## Practical Applications and Benefits

The graph-based knowledge representation allows for complex queries requiring holistic data understanding, such as:
- "Summarize the auth system of my app"
- "Refactor a component that spans multiple files"

Benefits of this approach include:
1. Contextual understanding of data relationships
2. Efficient information retrieval
3. Holistic analysis of data structure
4. Flexible querying across the knowledge base
5. Decentralized, resilient backend architecture
6. Real-time updates and progressive responses

This graph-based approach, enhanced by GraphRAG and Nostr communication, forms the core of OpenAgents' ability to provide intelligent, context-aware responses to user queries and commands.