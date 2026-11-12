/*
1. Define a load balancer server which asks to listen on a port using net.Listen()
2. Start accepting a request from a client and verify if the request is reaching my server
3. Introduce 1 backend server and connect that server to the load balancer
4. Create a bookkeeping mechanism for the load balancer? (or pass through.. this is the main plumbing of the entire project)
5. Set up communication between both the connections
6. Write code for closing connections (all situations we covered)
7. Introduce a second server.. same mechanism.. but now we introduce actual load balancing..
8. Later add more features... but this is the baseline
*/

/*
LB LEVEL 2: scenarios this code should handle
=============================================
Each line is a claim. For each one, write how you'll prove it before you
claim it.

DISTRIBUTION
  - Clients that arrive one at a time spread across every backend.
  - No backend gets picked twice in a row while another one is sitting unused.
  - The split stays fair when clients arrive at the same instant.
  - A backend still busy with long-lived clients isn't handed more work than an idle one.
  - Picking a backend stays cheap whether there are 3 backends or 3,000.

CORRECTNESS UNDER CONCURRENCY
  - Two clients connecting at the same moment never corrupt the proxy's view of the world.
  - The proxy's bookkeeping matches reality after every client has left.
  - Nothing the proxy counts drifts over a million connections.

FAILURE
  - One backend being down doesn't take the proxy down.
  - A client never hangs forever because of a backend's problem.
  - A backend leaving mid-conversation affects only the clients it was serving.

LIFECYCLE
  - Every connection that opens is eventually closed, on both sides.
  - Every goroutine that starts eventually ends, and you can say why.
  - Adding a fourth backend is a config change, not a code change.
  - Swapping the balancing strategy doesn't touch the proxying code.

STILL OWED FROM LEVEL 1
  - Own buffer per direction: why?
  - Lines in `ss -tn` with all three connected: how many?
  - The backend's idea of who its client is: what?
  - One 10,000-byte line: how many chunks, and how many lines?
*/

package main

import "net"

func main() {
    port_listener, err := net.Listen("tcp", "127.0.0.01:8080")
}