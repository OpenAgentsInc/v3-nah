# GraphRAG Implementation for OpenAgents

_Moved over from previous codebase; outdated, we won't use Convex at least to start - will use SQLite_

## Introduction
GraphRAG is a structured, hierarchical approach to Retrieval Augmented Generation (RAG) that can significantly enhance the capabilities of our AI productivity dashboard. By implementing a basic version of GraphRAG, we can improve the system's ability to understand and reason about the complex relationships within our codebase and user data.

## How GraphRAG Can Help OpenAgents

1. **Enhanced Code Understanding**: By creating a knowledge graph of our codebase, we can help the AI better understand the relationships between different components, functions, and data structures.

2. **Improved Query Responses**: GraphRAG can provide more accurate and contextually relevant responses to user queries about the codebase, features, or documentation.

3. **Holistic View of User Data**: By applying GraphRAG to user threads and messages, we can create a more comprehensive understanding of user interactions and conversations over time.

4. **Better Thread Management**: GraphRAG can help in generating more accurate thread titles and summaries by considering the broader context of conversations.

5. **Advanced Search Capabilities**: Implementing GraphRAG can enable more sophisticated search functionality across threads and messages.

## Key Concepts

1. **Files**: Individual source code files within a repository. They are the basic units of our codebase and contain the actual code content.

2. **Entities**: Distinct elements within the code, such as functions, classes, variables, or modules. Entities are the building blocks of our knowledge graph.

3. **Relationships**: Connections between entities, representing how different parts of the code interact or depend on each other. Examples include function calls, inheritance, or imports.

4. **Communities**: Groups of related entities that form a cohesive unit or serve a common purpose within the codebase. Communities help in understanding the high-level structure and organization of the code.

The importance of these concepts lies in their ability to capture the structure and semantics of the codebase, enabling more intelligent and context-aware interactions with the AI system.

## Implementation Plan

1. **Define TextUnits**:
   - Use individual GitHub files in our repo as TextUnits for the codebase.
   - For user data, consider using individual messages or entire threads as TextUnits.

2. **Entity Extraction**:
   - For codebase: Extract functions, classes, variables, and file names as entities.
   - For user data: Extract key concepts, user names, and important terms from messages and threads.

3. **Relationship Extraction**:
   - For codebase: Identify imports, function calls, and data flow between components.
   - For user data: Track conversation flow, topic relationships, and user interactions.

4. **Knowledge Graph Construction**:
   - Use Convex's existing capabilities, including vector similarity search, to store and query entities and relationships.
   - Implement a basic version of the Leiden algorithm for community detection, if possible within Convex's constraints.

5. **Summary Generation**:
   - Develop a system to generate summaries for code components and user threads.
   - Use existing AI capabilities to create concise descriptions of entities and communities.

6. **Query System Enhancement**:
   - Modify the existing chat system to incorporate GraphRAG-based retrieval.
   - Implement both global and local search functionalities as described in the GraphRAG documentation.

7. **Integration with Existing Systems**:
   - Update the `useChat.ts` hook to incorporate GraphRAG-enhanced responses.
   - Modify relevant Convex functions (e.g., `saveChatMessage.ts`, `getThreadMessages.ts`) to work with the new graph structure.

8. **UI Updates**:
   - Enhance the chat pane to visualize relationships between messages and concepts.
   - Add a new pane for exploring the knowledge graph of the codebase.

## Convex Schemas for GraphRAG Implementation

The following schemas have been added to the Convex configuration to support GraphRAG:

```typescript
graphRagCommunities: defineTable({
  name: v.string(),
  description: v.string(),
  createdAt: v.number(),
  updatedAt: v.number(),
  members: v.array(v.string()),
  repositories: v.array(v.string()),
}),

graphRagNodes: defineTable({
  communityId: v.id("graphRagCommunities"),
  type: v.string(),
  name: v.string(),
  content: v.string(),
  embedding: v.array(v.number()),
  metadata: v.object({
    repository: v.string(),
    filepath: v.string(),
    lineStart: v.optional(v.number()),
    lineEnd: v.optional(v.number()),
  }),
}),

graphRagEdges: defineTable({
  communityId: v.id("graphRagCommunities"),
  sourceId: v.id("graphRagNodes"),
  targetId: v.id("graphRagNodes"),
  type: v.string(),
  weight: v.optional(v.number()),
}),
```

These schemas allow us to:

1. Store GraphRAG communities with their associated metadata.
2. Represent nodes in the knowledge graph, including their content, type, and embedding.
3. Capture relationships between nodes as edges, with types and optional weights.
4. Associate nodes and edges with specific communities for organized knowledge representation.

## Next Steps

1. Implement functions to populate the GraphRAG tables based on the existing codebase and user data.
2. Develop algorithms for entity extraction and relationship identification specific to our codebase structure.
3. Create a system for generating and updating embeddings for nodes in the knowledge graph.
4. Implement community detection algorithms to group related nodes.
5. Modify the existing chat system to utilize the GraphRAG structure for more context-aware responses.
6. Develop visualization tools for exploring the knowledge graph.
7. Create a system for keeping the GraphRAG data up-to-date as the codebase evolves.
8. Implement advanced query mechanisms that leverage the graph structure for more accurate information retrieval.
9. Integrate the GraphRAG system with the existing credit and user management systems.
10. Develop comprehensive testing strategies for the GraphRAG implementation.

By implementing GraphRAG using these schemas and following the outlined steps, we can significantly enhance OpenAgents' ability to understand and reason about codebases and user interactions. This approach will provide users with more insightful and context-aware responses, improving the overall functionality and user experience of the AI productivity dashboard.
