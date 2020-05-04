package api ;

import java.io.* ;
import java.util.* ;
import java.net.* ;

import org.json.* ;
import org.restlet.resource.*;
import org.restlet.representation.* ;
import org.restlet.ext.jackson.* ;
import org.restlet.ext.json.* ;
import org.restlet.data.* ;
import org.restlet.* ;
import java.util.StringTokenizer;

import java.util.concurrent.ConcurrentLinkedQueue ;


public class Sync implements Runnable {

	private AdminServer server = AdminServer.getInstance() ;
	private ConcurrentLinkedQueue<SyncRequest> sync_queue ;
	private String sync_node ;

	public Sync( String node, ConcurrentLinkedQueue<SyncRequest> queue ) {
		sync_node = node ;
		sync_queue = queue ;	
	}

    // Background Thread
	@Override
	public void run() {

		SyncRequest syncObject = null ;

		while (true) {
			try {
				// sleep for 5 seconds
				try { Thread.sleep( 5000 ) ; } catch ( Exception e ) {}  

				//System.out.println( "CHECK SYNC QUEUE: " + sync_node + "..." ) ;

				if ( !sync_queue.isEmpty() ) {
					byte[] addr = ipAddrConverter(sync_node);
					InetAddress inet = InetAddress.getByAddress(addr) ;
					if (inet.isReachable(1000))
					{
					// check sync queue for work
						syncObject = sync_queue.peek() ;	

						// try to sync to peer node...
						System.out.println ( 	  "SYNC: " + "[" + server.getMyHostname() + "]" + " -> " + sync_node  
												+ " Document Key: " + syncObject.key 
												+ " vClock: " + Arrays.toString(syncObject.vclock) ) ;

						ClientResource client = server.getSyncClient( sync_node ) ;
						client.post( new JacksonRepresentation<SyncRequest>(syncObject), MediaType.APPLICATION_JSON);

						// remove head of queue if successfull
						syncObject = sync_queue.poll() ;
					}
					else
					{
					// if sync error, leave in queue for retry
						System.out.println ( 	  "SYNC: Node: " + sync_node + " is not reachable, RETRY later" ) ;
					}
				}				

			} catch (Exception e) {
				e.printStackTrace(); 
				System.out.println( e ) ;
			}			
		}
	}
	
	public byte[] ipAddrConverter(String ipAddress){
		byte[] result = new byte[4];
		StringTokenizer tokens = new StringTokenizer(ipAddress,".");
		int intAddress = 0;
		while(tokens.hasMoreTokens()){
			String octect  =  tokens.nextToken();
			int octectValue = Integer.valueOf(octect).intValue();
			intAddress = (intAddress << 8) + octectValue;
		}
		result[3] = (byte) (intAddress &  0xFF);
		result[2] = (byte) ((intAddress >> 8) &  0xFF);
		result[1] = (byte) ((intAddress >> 16) &  0xFF);
		result[0] = (byte) ((intAddress >> 24)&  0xFF);
		return result;
	}

}


