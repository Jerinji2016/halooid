# Architecture Overview

The Halooid platform is designed with a modern, scalable architecture that enables efficient development and maintenance of multiple products within a single ecosystem. This document provides a high-level overview of the platform's architecture.

## System Architecture

The Halooid platform follows a microservices architecture pattern, with clear separation between different products while sharing common infrastructure and libraries.

```mermaid
graph TD
    Client[Client Applications] --> API[API Gateway]
    API --> Auth[Authentication Service]
    API --> TS[Taskake Service]
    API --> QS[Qultrix Service]
    API --> AS[AdminHub Service]
    API --> CS[CustomerConnect Service]
    API --> IS[Invantray Service]
    
    TS --> DB[(PostgreSQL)]
    QS --> DB
    AS --> DB
    CS --> DB
    IS --> DB
    
    TS --> Cache[(Redis)]
    QS --> Cache
    AS --> Cache
    CS --> Cache
    IS --> Cache
    
    subgraph "Backend Services"
        Auth
        TS
        QS
        AS
        CS
        IS
    end
    
    subgraph "Data Layer"
        DB
        Cache
    end
```

## Key Components

### API Gateway

The API Gateway serves as the entry point for all client applications, handling request routing, composition, and protocol translation. It provides a unified interface for clients to interact with the various backend services.

### Authentication Service

The Authentication Service manages user authentication and authorization across all products. It implements JWT-based authentication with role-based access control to ensure secure access to resources.

### Product Services

Each product (Taskake, Qultrix, AdminHub, CustomerConnect, and Invantray) has its own dedicated service that implements the business logic specific to that product. These services are implemented in Go and expose both REST and gRPC APIs.

### Data Layer

The data layer consists of PostgreSQL for persistent storage and Redis for caching. PostgreSQL stores all application data with proper schema design for each product, while Redis improves performance by caching frequently accessed data.

## Communication Patterns

### External Communication

External clients (web and mobile applications) communicate with the backend services through REST APIs, which are documented using OpenAPI specifications. This provides a standardized, easy-to-use interface for client applications.

### Internal Communication

For internal communication between services, gRPC is used for its efficiency and strong typing. This enables high-performance, low-latency communication between backend services.

```mermaid
sequenceDiagram
    participant Client
    participant Gateway as API Gateway
    participant Auth as Authentication Service
    participant Service as Product Service
    participant DB as Database
    
    Client->>Gateway: HTTP Request
    Gateway->>Auth: Validate Token (gRPC)
    Auth-->>Gateway: Token Valid
    Gateway->>Service: Forward Request (gRPC)
    Service->>DB: Query Data
    DB-->>Service: Return Data
    Service-->>Gateway: Response
    Gateway-->>Client: HTTP Response
```

## Frontend Architecture

### Web Frontend

The web frontend is built using Svelte with SvelteKit, providing a modern, reactive user interface. Each product has its own SvelteKit application, but they share common UI components, stores, and utilities.

```mermaid
graph TD
    SvelteKit[SvelteKit] --> Router[Router]
    SvelteKit --> Store[Svelte Stores]
    SvelteKit --> Components[UI Components]
    
    Router --> Pages[Pages]
    Store --> State[Application State]
    Components --> SharedUI[Shared UI Library]
    
    Pages --> ProductPages[Product-Specific Pages]
    State --> ProductState[Product-Specific State]
    SharedUI --> ProductUI[Product-Specific UI]
```

### Mobile Frontend

The mobile applications are built using Flutter, providing a native-like experience on both iOS and Android. The mobile apps follow a similar architecture to the web frontend, with shared components and services.

```mermaid
graph TD
    Flutter[Flutter] --> Widgets[Widgets]
    Flutter --> Providers[Providers/Riverpod]
    Flutter --> Services[Services]
    
    Widgets --> SharedWidgets[Shared Widgets]
    Providers --> AppState[Application State]
    Services --> APIServices[API Services]
    
    SharedWidgets --> ProductWidgets[Product-Specific Widgets]
    AppState --> ProductState[Product-Specific State]
    APIServices --> ProductAPI[Product-Specific API]
```

## Database Design

The database is designed to support all products while maintaining proper separation of concerns. Each product has its own schema, but common entities like users and organizations are shared across products.

```mermaid
erDiagram
    USERS ||--o{ USER_ROLES : has
    USERS ||--o{ ORGANIZATIONS : belongs_to
    ORGANIZATIONS ||--o{ PROJECTS : has
    
    PROJECTS ||--o{ TASKS : contains
    TASKS ||--o{ TASK_ASSIGNMENTS : has
    USERS ||--o{ TASK_ASSIGNMENTS : assigned_to
    
    ORGANIZATIONS ||--o{ EMPLOYEES : employs
    USERS ||--o{ EMPLOYEES : is
    
    ORGANIZATIONS ||--o{ CUSTOMERS : has
    USERS ||--o{ CUSTOMERS : is
    
    ORGANIZATIONS ||--o{ INVENTORY_ITEMS : owns
    INVENTORY_ITEMS ||--o{ INVENTORY_TRANSACTIONS : has
```

## Security Architecture

Security is a fundamental aspect of the Halooid platform. The platform implements multiple layers of security:

1. **Authentication**: JWT-based authentication for all API requests
2. **Authorization**: Role-based access control for fine-grained permissions
3. **Data Encryption**: Encryption of sensitive data at rest and in transit
4. **Input Validation**: Thorough validation of all user inputs
5. **Rate Limiting**: Protection against abuse and DoS attacks
6. **Audit Logging**: Comprehensive logging of security-relevant events

## Deployment Architecture

The Halooid platform is designed to be deployed in a containerized environment using Docker and Kubernetes, enabling scalability and resilience.

```mermaid
graph TD
    Internet[Internet] --> LB[Load Balancer]
    LB --> API[API Gateway]
    
    API --> AuthPod[Auth Service Pods]
    API --> TaskakePod[Taskake Service Pods]
    API --> QultrixPod[Qultrix Service Pods]
    API --> AdminPod[AdminHub Service Pods]
    API --> CRMPod[CustomerConnect Service Pods]
    API --> InvPod[Invantray Service Pods]
    
    AuthPod --> DB[(PostgreSQL Cluster)]
    TaskakePod --> DB
    QultrixPod --> DB
    AdminPod --> DB
    CRMPod --> DB
    InvPod --> DB
    
    AuthPod --> Cache[(Redis Cluster)]
    TaskakePod --> Cache
    QultrixPod --> Cache
    AdminPod --> Cache
    CRMPod --> Cache
    InvPod --> Cache
    
    subgraph "Kubernetes Cluster"
        API
        AuthPod
        TaskakePod
        QultrixPod
        AdminPod
        CRMPod
        InvPod
    end
```

## Next Steps

For more detailed information about specific aspects of the architecture, please refer to the following pages:

- [Backend Architecture](backend.md)
- [Frontend Architecture](frontend.md)
- [Mobile Architecture](mobile.md)
- [Database Design](database.md)
- [API Architecture](api.md)
