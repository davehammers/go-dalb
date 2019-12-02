# go-dalb
## Go Dynamic Application Load Balancer

This Application Load Balancer accepts incomming HTTP requests and distributes them to worker nodes.

Each worker node has an IP address, port and the maximum number of outstanding transactins that can be queued. The system monitors the work nodes performance by measuring the time taken for each transaction. 

The rebalancer(TBD) runs periodically and examines the performance of each worker node. The lower performing worker nodes will be given less data traffic while the better performing worker nodes will get an increase.

The scheduling algorithm uses a Weighted Round Robin calendar to avoid the computational overhead for determining the next available worker node using mathmatical formula.

By using a calendar WRR, the system behavior is deterministic inbetween rebalance periods.

(TODO) add system and worker node performance history tracking.

## Build
Go version go1.13.4 was used to build this application.

To build go-dalb, enter:

```sh
make
```
on the command line. The system will build and test each component.

The resulting binary is:

```sh
go-dalb
```
## Running
go-dalb creates 2 HTTP servers, each listening on different ports.

- data path - default port is 8080
- control path - default port is 8081

The different ports enable go-dalb management to be separate from the data forwarding path. The URL's used to manage go-dalb do not conflict with any URL that may occur in the load balancing path.

### Control path URL's
The following URL's are available to gather statistics and add worker nodes to go-dalb

GET		/scheduler	returns global scheduler statistics

GET		/node			returns node statistics for all worker nodes

POST	/node			Adds a worker node to the scheduler 

