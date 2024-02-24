# go-chord

DataServer Struct:
- OverlayNode {IP: string, Hash: int, PredecessorIP: string, FingerTable: []string, potentially RingKey: string, }
- DataNode

Main: 
    Run a gRPC client to be able to send requests to peers
    Join an existing ring or Create a new one.
    Initialize new server struct.
    Run a gRPC server to handle peer requests
    Run a REST API server to handle user requests
    Run a set of maintenance routines as go routines

    Clean up when all routines are done and exit

Potential Organizations:

### Data, Lookup?, Maintenance, etc.

Folder 1: Data
- Implement KV Handlers
- Implement Key Transfer RPC (Not sure where to put this really)
- Implement 

### User, Client, Server

Folder 1: User
- Implement KV Handlers
- Maybe even implement a method that encapsulates starting the external facing REST server for the overall server.

Folder 2: Client
- Implement maintenance routines using the gRPC client functionality.
- Implement something to encapsulate starting the client and connecting to the necessary nodes.

Folder 3: Server

- Implement the Remote Procedures defined in the proto file.
- Implement a method to encapsulate starting the gRPC server.

However, the KV routes are supposedly being done **as** a user by analogy to the other packages, so this organization is a bit weird. Something closer to Data and Overlay might be sufficient. In fact, the final server could be defined as a conglomeration of two independent structs with receiver functions defined in each folder - no, that wouldn't work as gRPC also needs access to data for data transfer.

**gRPC services need to be defined as functions of a data structure that contains both Data and Overlay information. REST services only need to be defined as functions of the Data information.**

New Struct:

ChordNode: {
    * DataServer {
        KVMap,
        KeyLookupStructure (BST, RBT, etc.)
    },
    IP: string, 
    Hash: int, 
    Capacity: int,
    PredecessorIP: string, 
    FingerTable: []string, 
    potentially RingKey: string
}

The server contains the overlay information directly (without abstraction) and the data information via the DataServer pointer.

The DataServer has Data Service Handlers attached (for user-facing solution) + the Chord Node implements gRPC Services and performs maintenance.

### Data, Overlay

The Data package defines the Data Server struct and then all the handlers + a method that starts listening to the specified port to handle user requests.

The Overlay package defines the Chord Node struct itself (imports DataServer from Data package), and all gRPC services, clients, maintenance, join/create.

The main package then just runs a main method which creates a ChordNode, joins/creates, and starts the DataServer.