![SuBFraMe - Logo](https://www.fuchstim.de/subframe/res/img/logo/logo-1-text_small.png#)
# SuBFraMe - Protocol

## Table of Contents
1. [Sending](#sending)
2. [Receiving](#receiving)
3. [StorageNode](#storagenode)
4. [CoordinatorNode](#coordinatornode)
5. [Bootstrapping](#bootstrapping)

## Client
### Sending
#### 1. Composing
To send a message, certain components are required.

`message: "Hello World"`        -> (optional) Text Content of the message

`attachment: "base64: (...)"`   -> (optional) Any attachments

(Either message, attachment or both are required)

#### 2. Sealing
The Message's contents are now placed inside an envelope

`envelope: { message: "Hello World", attachment: "base64: (...)", checksum: "(e.g. sha256sum or similar, of message + attachment)" }`

The envelope is now encrypted with the sender's private key

`envelope: { <encrypted gibberish> }`

and with the recipient's public key

`envelope: { <even more gibberish> }`

(This can be made faster, especially for large messages, with encrypting the message + attachments using a random passphrase and encrypting that with PK-Crypto)


The envelope is now assigned a unique ID:

`envelope: { id: "<recipient's ID (e.g. pubkey) + checksum of checksum>", content: <encrypted> }`

The envelope is now ready to be transmitted


#### 3. Transmission
The client now pushes the envelope to one or more StorageNodes:

`POST { url: "https://node-address/storage/put/<envelope-id>", body: "<envelope-content>" }`

After successfully receiving and storing the message, the StorageNode(s) announce to at least 3 random CoordinatorNodes that they know of and serve the message:

`GET { url: "https://node-address/coordinator/announce/<envelope-id>/<own-address>"}`

The CoordinatorNodes respond with either `"true"` or `"false"` (or an error). Depending on the result the StorageNode pushes the envelope to 1+ more StorageNode(s). This cycle repeats until the CoordinatorNetwork responds with `"false"`


### Receiving
#### 1. Transmission
The recipient queries any CoordinatorNode for new Messages for it's ID.

`GET { url: "https://node-address/coordinator/get/<client-id>"}`

CoordinatorNodes responds with any (unreceived) messages:

`[{ id: "<envelope1-id>", content: "<envelope1-content>", storageNodes: ["node1-address", "node2-address", "node3-address"]},(...)]`

The recipient can now query one (or multiple, for verification) of the listed StorageNodes for the message:

`GET { url: "https://node1-address/storage/get/<envelope1-id>" }`

to receive

`{ id: "<envelope1-id>", content: "<envelope1-content>"}`


#### 2. Decryption and Verification
The received envelope is now decrypted using the recipient's private key, and the sender's public key

`envelope: { message: "Hello World", attachment: "base64: (...)", checksum: "(e.g. sha256sum or similar, of message + attachment)" }`

The recipient now generates the checksum and checks whether it matches the one in the envelope, as well as whether the checksum of the checksum matches the one in the envelope-id. If they match, the message is displayed, else the recipient receives a warning.

The calculated checksum of the message is then transmitted to the CoordinatorNetwork (again - to one or more CoordinatorNodes for redundancy):

`GET { url: "https://coordinator-node/coordinator/verify/<envelope-id>/<checksum>" }`

If the CoordinatorNode is able to verify that the checksum of the received checksum matches the one in the envelope-id, the message is marked as received and removed from the CoordinatorNetwork's database.


#### 3. Deletion
Depending on the StorageNodes' settings, a message is deleted if either

1. It exceeds the maximum duration a StorageNode stores a Message (configurable)
    \- or -
2. It is marked as received in, or disappeared from, the CoordinatorNetwork's database

(`GET { url: "https://node-address/coordinator/get-status/<envelope-id>" }` returns message status (`0: not in database, 1: in database, -1: error`))


The latter is checked periodically, the minimum duration between checks is also configurable.

### StorageNode
A StorageNodes serves as file storage space for messages. It can receive and store, as well as serve messages.
It exposes a very basic set of endpoints:

#### `/storage/`
- `GET /storage/get/<id>`: Returns envelope, if present
- `POST /storage/put/<id> | body: <content>`: Stores message to node, if possible

#### `/control/`
- `GET /control/export-coordinator-nodes` and `GET /control/export-storage-nodes`: Exports known CoordinatorNodes and StorageNodes respectively (for bootstrapping new node)

### CoordinatorNode
A CoordinatorNode is part of the CoordinatorNetwork. This network holds a synchronous database with all current (not yet received) messages present in the network. To make this synchronization possible, the network is limited in size (max. ~ 20 Nodes?). 


Further technical details on how the database is kept in sync will follow


Nodes in the CoordinatorNetwork are dynamic, depending on their uptime and connection quality they may be kicked from the CoordinatorNetwork or leave intentionally, new Nodes can join the CoordinatorNetwork if it's current size allows for it. This impersistance in CoordinatorNetwork structure enhances the Network's security, but also creates the possibility of 'Lost' Nodes.

A 'Lost' Node does not know _any_ active CoordinatorNodes. If a Node encounters one or more inactive CoordinatorNodes in it's database, it is able to receive the updated, full list of CoordinatorNodes from one of the remaining, known and active CNodes. If the Node was unable to update the local list of CoordinatorNodes (e.g. after a longer downtime after which all previous CNs are now inactive), it would need to query one or more known StorageNodes for a recent list of active CoordinatorNodes. If this fails as well, the Node would need to be bootstrapped again.


A CoordinatorNode exposes a set of endpoints:

#### `/coordinator/`
- `GET /coordinator/get/<id>`: Returns list of StorageNodes holding Message with ID
- `GET /coordinator/verify/<id>/<verification-code>`: Verifies Message Reception
- `GET /coordinator/announce/<id>/<StorageNode-Address>`: Adds storageNode as server for message

#### `/control/`
- `GET /control/export-coordinator-nodes` and `GET /control/export-storage-nodes`: Exports known CoordinatorNodes and StorageNodes (for bootstrapping new member)


### Bootstrapping
To bootstrap a new client, it needs to be provided a ´bootstrap-node´. This can be any Node on the network.
This node now exports it's list of StorageNodes and CoordinatorNodes, the new Node writes it to it's database.
The Node can now decide to further query this new list of nodes to expand and update it, or to not do that.


Using specific ´bootstrap-nodes´, it is possible to run multiple SuBFraMe network simultaneously. If you were to carefully bootstrap Nodes with a very select number of nodes, you are theoretically able to hermetically isolate one SuBFraMe Network from another. As soon as just one Node on one network logs one Node from the other in it's database, however, the two networks merge.
