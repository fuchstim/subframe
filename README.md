[![fuchstim.de - Logo](https://fcdn.ftim.eu/res/img/logo_medium.png)](https://www.fuchstim.de)![SuBFraMe - Logo](https://www.fuchstim.de/subframe/res/img/logo/logo-1-nobg_medium.png)
# SuBFraMe Messaging Framework

1. [What is SuBFraMe?](#what-is-subframe)
2. [Download and Install](#download-and-install)
3. [Contributing](#contributing)
4. [How does it work?](#how-does-it-work)
5. [Sending a Message](#sending-a-message)
6. [Advantages and Disadvantages](#advantages)

#### What is SubFraMe?
SuBFraMe stands for **S**ec**u**rely **B**roadcasted and **Fra**gmented **Me**ssaging. It is a concept for an entirely decentralized network used primarily for email-like communication and static file-sharing.

#### Download and Install
Currently, there is no working Version of SuBFraMe available. If you wish to contribute to make SuBFraMe a reality, see [#Contributing](#contributing)

#### Contributing
If you are interested in the SuBFraMe Project and want to contribute, you can find a guide to get you set up in the [CONTRIBUTING.md](CONTRIBUTING.md) file.
Any help is greatly appreciated and will help make communication more secure and uncensored!

#### How does it work?
The SubFraMe Network consists of two kinds of Nodes: StorageNodes and CoordinatorNodes (SN and CN respectively). CNs represent a small fragment of and are therefore part of the SN-Network. 
##### StorageNodes
When a message is to be sent, it is transmitted to a StorageNode. This Node will store it to it's local disk for the message file to be served on request. It will also log the unique MessageID to it's local database and announce to the Coordinator-Network that it has the Message File in question. It then proceeds to transmit the Message to a number of other StorageNodes, for them to also store the Message and announce that they are able to serve it to the recipient. 
A StorageNode periodically tries to register as a CoordinatorNode, to fill any open positions in the CoordinatorNetwork. A successful registration will sync the CoordinatorNetwork Database to the StorageNode.
StorageNodes with a good connection will be prioritized on registration as CoordinatorNode, as well as being selected as StorageNode for a Message.
##### CoordinatorNodes
The Coordinator-Network represents the 'Brain' of the SubFraMe Network, and is strictly limited in size. It possesses a synchronous database keeping track of known Messages and their StorageNodes. That way, when a recipient requests a list of new Messages, it can query any CoordinatorNode and get a list of relevant Messages and their Storage Locations.
CoordinatorNodes with a bad connection may be kicked from the CoordinatorNetwork, yet still function as a StorageNode.

#### Sending a Message
##### Client-App
1. Compose Message and add Attachments
2. Create Envelope with Message Body, Attachments, Checksum, Confimation Key and other relevant metadata
3. Encrypt the Envelope and generate a unique MessageID from the RecipientID and a Hash of the Messages Confirmation Key
4. Transmit the MessageFile to any known StorageNode

##### StorageNode
1. Receive the Message File
2. Store Message to local disk
3. Log MessageID to local database, broadcast own ID as server for Message File with MessageID to CoordinatorNode
4. Transmit MessageFile to defined number of other StorageNodes

##### CoordinatorNode
1. Receive ID of StorageNode registering as MessageFileServer
2. Log StorageNodeID and MessageID to local database
3. Synchronize database with all other CoordinatorNodes

##### Recipient
1. Request new Messages with MessageID containing ID of Recipient
2. Receive List of new Messages and their StorageNodes
3. Download Message File from one or more of the relevant StorageNodes
4. Decrypt Message, verify checksum
5. Verify that Confimation Key in envelope matches Hash in MessageID
6. Announce to CoordinatorNode that Message has been received and decrypted, transmit Confirmation Key
7. Display Message Body and Attachments

##### CoordinatorNode
1. If Confirmation Key Hash matches MessageID, mark Message as received, verified and decrypted

##### StorageNode
1. Periodically check Message Status against CoordinatorNetwork, remove file from disk if Message has been received and decrypted or if MessageTTL has passed

#### Schematic: Sending a message
[![Schematic - Sending](https://www.fuchstim.de/subframe/res/img/schematic/sending_preview.png)](https://www.fuchstim.de/subframe/res/img/schematic/sending.png) (Click to enlarge)

#### Schematic: Receiving a message
[![Schematic - Receiving](https://www.fuchstim.de/subframe/res/img/schematic/receiving_preview.png)](https://www.fuchstim.de/subframe/res/img/schematic/receiving.png) (Click to enlarge)

#### Advantages
- Network is entirely decentralized and pretty much impossible to censor or block
- SuBFraMe allows for secure communication through Public-Key-Encryption and various checksum verifications
- Message Flow is hard to trace, as sender is only identifiable by recipient after decrypting the envelope
- ...

#### Disadvantages
-  Message Propagation can be rather slow. making real-time-communication not possible
-  Pushing messages to the recipients device is not possible, it will have to query the CoordinatorNetwork for new Messages 
-  Sender's identify cannot be reliably verified without a secure way of exchanging PublicKeys; meeting physically and using uncompromised systems
-  Using multiple devices to receive messages for the same address will be rather complicated, as the message is removed from the Network once it has been decrypted and verified, or reached its TTL
-  ...
