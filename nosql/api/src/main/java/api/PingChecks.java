package api ;

import java.io.* ;
import java.util.* ;
import java.net.*;

public class PingChecks implements Runnable {

	private AdminServer server = AdminServer.getInstance() ;

    // Background Thread
	@Override
	public void run() {
		while (true) {
			try {
				// sleep for 1 second
				try { Thread.sleep( 1000 ) ; } catch ( Exception e ) {}  
				// ping & sync nodes
				Collection<Node> nodes = server.getNodes() ;
    			for (Iterator<Node> iterator = nodes.iterator(); iterator.hasNext();) {
    				Node n = iterator.next() ;
    				String my_host = server.getMyHostname() ;
    				if ( !n.id.equals(my_host) ) {
						
	    				try {
							byte[] addr = ipAddrConverter(n.name);
	    					InetAddress inet = InetAddress.getByAddress(addr) ;
	    					if (inet.isReachable(1000)) {
								System.out.println( "Ping Node [id:" + n.id + "  name:" + n.name + "] ==> Node Up!" ) ;
								server.nodeUp( n.id ) ;      
							} else {
								System.out.println( "Ping Node [id:" + n.id + "  name:" + n.name + "] ==> Node Down!" ) ;
								server.nodeDown( n.id ) ;
							}
	  							
	  					} catch (Exception e) {
	  						server.nodeDown( n.id ) ; 
	  						System.out.println( e ) ;
						}   		    					
    				} else {
    					server.nodeSelf( n.id ) ;
    				}
    				
  				} 
			} catch (Exception e) {
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



