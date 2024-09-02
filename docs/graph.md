# Graphs in OpenAgents

The main screen of the mobile app displays a graph representing the agent's knowledge base. This graph structure is fundamental to how OpenAgents processes and retrieves information.

## Graph Theory Approach

The knowledge base is modeled using graph theory concepts, with information represented as nodes (vertices) connected by edges. This structure creates a rich, interconnected representation of data.

## Representing Data as a Graph

Any corpus of data, such as a codebase, can be represented as a graph:

- Nodes represent discrete pieces of information or concepts
- Edges represent relationships between these nodes
- Nodes may be grouped into communities or topics
- Each community or topic can have its own summary

## Vector Embeddings and Similarity Search

Each node in the graph is associated with a vector embedding, enabling:

- Cosine similarity search between nodes
- Specific information retrieval

This approach, inspired by Microsoft's GraphRAG paper, combines graph-based representation with vector-based similarity search.

## GraphRAG Implementation

OpenAgents implements GraphRAG (Graph-based Retrieval Augmented Generation), a structured, hierarchical approach that improves understanding and reasoning about complex relationships within codebases and user data.

Key aspects of the GraphRAG implementation:

1. **Code Understanding**: Creates a knowledge graph of the codebase to map relationships between components, functions, and data structures.
2. **Query Responses**: Provides contextually relevant responses to user queries about the codebase, features, or documentation.
3. **User Data Analysis**: Applies GraphRAG to user threads and messages for comprehensive understanding of user interactions over time.
4. **Advanced Search**: Enables sophisticated search functionality across threads and messages.

The implementation process includes:
- Defining text units
- Extracting entities and relationships
- Constructing a knowledge graph
- Generating summaries
- Enhancing the query system to leverage the graph structure

## Nostr-based Communication

OpenAgents uses Nostr for communication between the mobile app (client) and the relay (distributed backend). This approach leverages Nostr's decentralized nature and event-based communication model.

### Event Types

OpenAgents uses custom Nostr event kinds for various operations:

1. **Agent Command Request (kind 5838)**: When a user requests the agent to perform an action.
2. **Agent Command Response (kind 6838)**: Updates and responses from the agent as it processes the request.
3. **Speech-to-Text Request (kind 5252)**: For voice command transcription.
4. **Speech-to-Text Response (kind 6252)**: Transcription results.

### Communication Flow

1. The user initiates an action in the mobile app (e.g., a voice command).
2. The app sends a kind 5252 event to connected relay(s) for speech-to-text conversion.
3. The relay processes the audio and responds with a kind 6252 event containing the transcription.
4. The app then sends a kind 5838 event with the transcribed command to the relay.
5. The relay begins processing the command, which may involve graph-building (indexing) to gather necessary context.
6. As the relay processes the request, it sends one or more kind 6838 events back to the app with updates or partial results.
7. Once processing is complete, the relay sends a final kind 6838 event with the full response.

This event-based system allows for asynchronous, real-time communication between the app and the distributed backend, enabling responsive and flexible interactions.

## Practical Applications

This graph-based knowledge representation enables complex queries requiring holistic data understanding. Examples:

- "Summarize the auth system of my app"
- "Refactor a component that spans multiple files"

The agent uses the graph to find relevant context for building a prompt, which is then used with an LLM (Large Language Model) to generate answers or perform tasks.

## Benefits

1. Contextual Understanding: The graph structure reveals relationships between different parts of the data.
2. Efficient Retrieval: Vector embeddings and graph traversal enable quick information retrieval.
3. Holistic Analysis: Grouping nodes into communities allows higher-level understanding of data structure.
4. Flexible Querying: Users can ask questions that require synthesizing information from multiple parts of the knowledge base.
5. Decentralized Communication: Nostr-based messaging enables a distributed and resilient backend architecture.
6. Real-time Updates: The event-based system allows for immediate feedback and progressive responses to user queries.

This graph-based approach, enhanced by GraphRAG and Nostr communication, forms the core of OpenAgents' ability to provide intelligent, context-aware responses to user queries and commands, improving the functionality and user experience of the AI productivity dashboard.