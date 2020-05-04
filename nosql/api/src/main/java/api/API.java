package api ;

import nojava.* ;

import java.util.* ;
import java.io.* ;

import java.util.concurrent.BlockingQueue ;
import java.util.concurrent.LinkedBlockingQueue ;
import java.util.concurrent.ConcurrentHashMap ;
import java.util.Collection ;
import java.util.Random;
import java.util.UUID;
import java.util.concurrent.BlockingQueue;

import org.json.* ;
import org.restlet.resource.*;
import org.restlet.representation.* ;
import org.restlet.ext.json.* ;
import org.restlet.data.* ;


public class API implements Runnable {

	// queue of new documents
    private static BlockingQueue<Document> CREATE_QUEUE = new LinkedBlockingQueue<Document>() ;						

    // key to record map
    private static ConcurrentHashMap<String,Document> KEYMAP_CACHE = new ConcurrentHashMap<String,Document>() ;
    private static ConcurrentHashMap<String,Document> DELETEKEY_CACHE = new ConcurrentHashMap<String,Document>() ; 	

    // document updater map
    // private static HashMap<String, Integer> docUpdaterNode = new HashMap<String, Integer>() ;

    // Background Thread
	@Override
	public void run() {
		while (true) {
			try {
				// sleep for 5 seconds
				try { Thread.sleep( 5000 ) ; } catch ( Exception e ) {}  

				// process any new additions to database
				Document doc = CREATE_QUEUE.take();
  				SM db = SMFactory.getInstance() ;
        		SM.OID record_id  ;
        		String record_key ;
        		SM.Record record  ;
        		String jsonText = doc.json ;
            	int size = jsonText.getBytes().length ;
            	record = new SM.Record( size ) ;
            	record.setBytes( jsonText.getBytes() ) ;
            	record_id = db.store( record ) ;
            	record_key = new String(record_id.toBytes()) ;
            	doc.record = record_key ;
            	doc.json = "" ;
            	KEYMAP_CACHE.put( doc.key, doc ) ;    
                System.out.println( "Created Document: " + doc.key ) ;
                
                // sync nodes
                AdminServer.syncDocument( doc.key, "create" ) ; 

			} catch (InterruptedException ie) {
				ie.printStackTrace() ;
			} catch (Exception e) {
				System.out.println( e ) ;
			}			
		}
	}    


    public static Document[] get_hashmap() {
    	return (Document[]) KEYMAP_CACHE.values().toArray(new Document[0]) ;
    }


    public static void save_hashmap() {
		try {
		  FileOutputStream fos = new FileOutputStream("index.db");
		  ObjectOutputStream oos = new ObjectOutputStream(fos);
		  oos.writeObject(KEYMAP_CACHE);
		  oos.close();
		  fos.close();
		} catch(IOException ioe) {
		  ioe.printStackTrace();
		}
    }


    public static void load_hashmap() {
		 try {
		     FileInputStream fis = new FileInputStream("index.db") ;
		     ObjectInputStream ois = new ObjectInputStream(fis) ;
		     KEYMAP_CACHE = (ConcurrentHashMap) ois.readObject();
		     ois.close() ;
		     fis.close() ;
		  } catch(IOException ioe) {
		     ioe.printStackTrace() ;
		  } catch(ClassNotFoundException c) {
		     System.out.println("Class not found");
		     c.printStackTrace() ;
		  }
    }


   public static void sync_document(SyncRequest sync) throws DocumentException {

        String key = sync.key ;
        String value = sync.json ;
        String[] vclock = sync.vclock ;
        String command = sync.command ;

        try {

            AdminServer server = AdminServer.getInstance() ;
            String my_host = server.getMyHostname() ;
            int my_index = server.nodeIndex( my_host ) ;
            int senderIndex = server.nodeIndex( sync.vclock[0] ) ;
            char dChar, sChar ;
            System.out.println("API sync_document(): [ From Node: " + senderIndex + ", key: " + key + ", command: " + command + " ]");
            switch ( command ) {
                case "create":
                    System.out.println("API : Syncing Create request") ;
                    boolean createNeeded = true ;
                    Document doc = KEYMAP_CACHE.get( key ) ;
                    if ( doc != null && doc.record != null )
                    {   
                        for(int i = vclock.length-1 ; i > 0 ; --i)
                        {
                            System.out.println("API : Beginning of Update For Loop ") ;
                            dChar = (doc.vclock[i] == null) ? '0' : doc.vclock[i].charAt(doc.vclock[i].length()-1) ;
                            sChar = (vclock[i] == null) ? '0' : vclock[i].charAt(vclock[i].length()-1) ;
                            System.out.println("API : doc v[" + i + "] : " + doc.vclock[i] + ", sync v[" + i + "] : " + vclock[i]) ;
                            if ( (sChar < dChar) && (i > senderIndex))
                            {
                                System.out.println("API : ingoring Create Sync Request from Node: " + senderIndex + ", of Node: " + i + " as doc already present") ;
                                createNeeded = false ;
                                break;
                            }
                        }
                    }
                    else
                    {
                        doc = DELETEKEY_CACHE.get( key ) ;
                        if( doc != null && doc.vclock != null )
                        {
                            for(int i = vclock.length-1 ; i > 0 ; --i)
                            {
                                System.out.println("API : Beginning of Update For Loop ") ;
                                dChar = (doc.vclock[i] == null) ? '0' : doc.vclock[i].charAt(doc.vclock[i].length()-1) ;
                                sChar = (vclock[i] == null) ? '0' : vclock[i].charAt(vclock[i].length()-1) ;
                                System.out.println("API : doc v[" + i + "] : " + doc.vclock[i] + ", sync v[" + i + "] : " + vclock[i]) ;
                                if ( (sChar < dChar) && (i > senderIndex))
                                {
                                    System.out.println("API : ingoring Create Sync Request from Node: " + senderIndex + ", of Node: " + i + " as doc already deleted") ;
                                    createNeeded = false ;
                                    break;
                                }
                            }
                            if(createNeeded)
                            {
                                DELETEKEY_CACHE.remove( key );
                            }
                        }
                    }
                    if(createNeeded)
                    {
                        doc = new Document() ;
                        doc.vclock[0] = my_host ;
                        doc.vclock[1] = vclock[1] ;
                        doc.vclock[2] = vclock[2] ;
                        doc.vclock[3] = vclock[3] ;
                        doc.vclock[4] = vclock[4] ;
                        doc.vclock[5] = vclock[5] ;
                        doc.vclock[my_index] = my_host + ":" + Integer.toString(0) ;
                        SM db = SMFactory.getInstance() ;
                        SM.OID record_id  ;
                        SM.Record record  ;
                        String jsonText = value ;
                        int size = jsonText.getBytes().length ;
                        record = new SM.Record( size ) ;
                        record.setBytes( jsonText.getBytes() ) ;
                        record_id = db.store( record ) ;
                        String record_key = new String(record_id.toBytes()) ;
                        doc.record = record_key ;
                        doc.json = "" ;
                        doc.key = key ;
                        KEYMAP_CACHE.put( key, doc ) ;
                        System.out.println( "SYNC: Created Document Key: " + key 
                                        + " Record: " + record_key 
                                        + " vClock: " + Arrays.toString(doc.vclock) 
                                    ) ;
                    }
                    break ;
                case "update":
                    System.out.println("API : Syncing Update request") ;
                    boolean updateNeeded = true ;
                    doc = KEYMAP_CACHE.get( key ) ;
                    if ( doc != null && doc.record != null )
                    {   
                        for(int i = vclock.length-1 ; i > 0 ; --i)
                        {
                            System.out.println("API : Beginning of Update For Loop ") ;
                            dChar = (doc.vclock[i] == null) ? '0' : doc.vclock[i].charAt(doc.vclock[i].length()-1) ;
                            sChar = (vclock[i] == null) ? '0' : vclock[i].charAt(vclock[i].length()-1) ;
                            System.out.println("API : doc v[" + i + "] : " + doc.vclock[i] + ", sync v[" + i + "] : " + vclock[i]) ;
                            if ( (sChar < dChar) && (i > senderIndex) )
                            {
                                System.out.println("API : ingoring update from Node: " + senderIndex + ", of Node: " + i) ;
                                updateNeeded = false ;
                                break;
                            }
                        }
                    }
                    else
                    {
                        doc = DELETEKEY_CACHE.get( key ) ;
                        if( doc != null )
                        {
                            for(int i = vclock.length-1 ; i > 0 ; --i)
                            {
                                System.out.println("API : Beginning of Update For Loop ") ;
                                dChar = (doc.vclock[i] == null) ? '0' : doc.vclock[i].charAt(doc.vclock[i].length()-1) ;
                                sChar = (vclock[i] == null) ? '0' : vclock[i].charAt(vclock[i].length()-1) ;
                                System.out.println("API : doc v[" + i + "] : " + doc.vclock[i] + ", sync v[" + i + "] : " + vclock[i]) ;
                                if ( (sChar < dChar) && (i > senderIndex) )
                                {
                                    System.out.println("API : ingoring update from Node: " + senderIndex + ", of Node: " + i) ;
                                    updateNeeded = false ;
                                    break;
                                }
                            }
                            if(updateNeeded)
                            {
                                DELETEKEY_CACHE.remove( key ) ;
                                sync.command = new String("create") ;
                                sync_document(sync) ;
                                break;
                            }
                        }
                        else
                            break ;
                    }
                    
                    if(updateNeeded)
                    {
                        System.out.println("API : updating document from Node: " + senderIndex ) ;
                        String record_key = doc.record ;
                        System.out.println("API : getting DB instance") ;
                        SM db = SMFactory.getInstance() ;    
                        System.out.println("API : getting DB record_id") ;
                        SM.OID record_id = db.getOID( record_key.getBytes() ) ;
                        String jsonText = value ;
                        int size = jsonText.getBytes().length ;
                        try {
                            System.out.println("API : UPDATE creating new record") ;
                            SM.Record record = new SM.Record( size ) ;
                            record.setBytes( jsonText.getBytes() ) ;
                            SM.OID update_id = db.update( record_id, record ) ;
                            // docUpdaterNode.put(key, senderIndex);
                        }
                        catch (SM.NotFoundException nfe) {
                            System.out.println("API : Document Not Found: " + key) ;
                            throw new DocumentException( "Document Not Found: " + key ) ;
                        } 
                        catch (Exception e) {
                                System.out.println("API : Some exception occured") ;
                                throw new DocumentException( e.toString() ) ;           
                        }
                    }
                    System.out.println("API : merging vclock from Node: " + senderIndex ) ;
                    for(int i = vclock.length - 1 ; i > 0 ; --i )
                    {
                        dChar = (doc.vclock[i] == null) ? '0' : doc.vclock[i].charAt(doc.vclock[i].length()-1) ;
                        sChar = (vclock[i] == null) ? '0' : vclock[i].charAt(vclock[i].length()-1) ;
                        System.out.println("API : doc v[" + i + "] : " + doc.vclock[i] + ", sync v[" + i + "] : " + vclock[i]) ;
                        if( sChar > dChar )
                        {
                            doc.vclock[i] = vclock[i] ;
                        }
                    }
                    System.out.println( "SYNC: Updated Document Key: " + key 
                                    + " Record: " + doc.record
                                    + " vClock: " + Arrays.toString(doc.vclock) 
                                ) ;
                    break ;
                case "delete":
                    boolean deleteNeeded = true ;
                    System.out.println( "SYNC Delete Document: " + key ) ;
                    doc = KEYMAP_CACHE.get( key ) ;
                    if ( doc != null && doc.record != null )
                    {
                        for(int i = vclock.length-1 ; i > 0 ; --i)
                        {
                            System.out.println("API : Beginning of conflict For Loop in delete command ") ;
                            dChar = (doc.vclock[i] == null) ? '0' : doc.vclock[i].charAt(doc.vclock[i].length()-1) ;
                            sChar = (vclock[i] == null) ? '0' : vclock[i].charAt(vclock[i].length()-1) ;
                            System.out.println("API : doc v[" + i + "] : " + doc.vclock[i] + ", sync v[" + i + "] : " + vclock[i]) ;
                            if( (sChar < dChar) && ( i > senderIndex)) 
                            {
                                System.out.println("API : Due to conflict ignore delete from Node: " + senderIndex + ", of Node: " + i) ;
                                deleteNeeded = false;
                                break;
                            }
                        }
            
                    }
                    else
                    {
                        doc = DELETEKEY_CACHE.get( key ) ;
                        if( doc == null )
                        {
                            doc = new Document();
                            doc.key = key ;
                            doc.vclock[0] = my_host ;
                            doc.vclock[1] = vclock[1] ;
                            doc.vclock[2] = vclock[2] ;
                            doc.vclock[3] = vclock[3] ;
                            doc.vclock[4] = vclock[4] ;
                            doc.vclock[5] = vclock[5] ;
                            doc.json = null ;
                        }
                        else
                        {
                            for(int i = vclock.length - 1 ; i > 0 ; --i )
                            {
                                dChar = (doc.vclock[i] == null) ? '0' : doc.vclock[i].charAt(doc.vclock[i].length()-1) ;
                                sChar = (vclock[i] == null) ? '0' : vclock[i].charAt(vclock[i].length()-1) ;
                                System.out.println("API : doc v[" + i + "] : " + doc.vclock[i] + ", sync v[" + i + "] : " + vclock[i]) ;
                                if( sChar > dChar )
                                {
                                    doc.vclock[i] = vclock[i] ;
                                }
                            }
                        }
                        DELETEKEY_CACHE.put(key, doc) ;
                        deleteNeeded = false ;
                    }
                    if(deleteNeeded)
                    {
                        String record_key = doc.record ;
                        SM db = SMFactory.getInstance() ;    
                        SM.OID record_id = db.getOID( record_key.getBytes() ) ;
                        try {
                            db.delete( record_id ) ;
                            Document delDoc = new Document() ;
                            delDoc.key = (doc.key == null ) ? null : new String(doc.key) ;
                            delDoc.vclock[0] = my_host ;
                            delDoc.vclock[1] = (vclock[1] == null) ? null : new String(vclock[1]) ;
                            delDoc.vclock[2] = (vclock[2] == null) ? null : new String(vclock[2]) ;
                            delDoc.vclock[3] = (vclock[3] == null) ? null : new String(vclock[3]) ;
                            delDoc.vclock[4] = (vclock[4] == null) ? null : new String(vclock[4]) ;
                            delDoc.vclock[5] = (vclock[5] == null) ? null : new String(vclock[5]) ;
                            delDoc.vclock[my_index] = (doc.vclock[my_index] == null) ? null : new String(doc.vclock[my_index]) ;
                            delDoc.json = null ;
                            DELETEKEY_CACHE.put(key, delDoc) ;
                            KEYMAP_CACHE.remove( key ) ;
                            System.out.println( "SYNC Document Deleted: " + key ) ;
                        } catch (SM.NotFoundException nfe) {
                            throw new DocumentException( "Document Not Found: " + key ) ;
                        } catch (Exception e) {
                            throw new DocumentException( e.toString() ) ;            
                        }
                    }
                    break ;
            }      

        } catch (Exception e) {
            System.out.println("API : Exception thrown") ;
            throw new DocumentException( e.toString() ) ;
        }

    }


    public static void create_document(String key, String value) throws DocumentException {
    	try {
	    	System.out.println( "Create Document: Key = " + key + " Value = " + value ) ;
	    	Document doc = new Document() ;
	    	doc.key = key ;
            AdminServer server = AdminServer.getInstance() ;
            String my_host = server.getMyHostname() ;
            System.out.println( "My Host Name: " + my_host ) ;
            doc.vclock[0] = my_host ;
            String my_version = my_host + ":" + Integer.toString(1) ;
            int my_index = server.nodeIndex( my_host ) ;
            System.out.println( "Node Index: " + my_index ) ;
            doc.vclock[my_index] = my_version ;
	    	KEYMAP_CACHE.put( key, doc ) ;
	    	doc.json = value ;
	        CREATE_QUEUE.put( doc ) ; 
            // docUpdaterNode.put(key,my_index) ;
	    	System.out.println( "New Document Queued: " + key ) ;    		
	    } catch (Exception e) {
	    	throw new DocumentException( e.toString() ) ;
	    }

    }


    public static String get_document(String key) throws DocumentException {
    	System.out.println( "Get Document: " + key ) ;
    	Document doc = KEYMAP_CACHE.get( key ) ;
    	if ( doc == null || doc.record == null )
    		throw new DocumentException( "Document Not Found: " + key ) ;
    	String record_key = doc.record ;
    	SM db = SMFactory.getInstance() ;
    	SM.OID record_id ;
        SM.Record found ;
		record_id = db.getOID( record_key.getBytes() ) ;
        try {
            found = db.fetch( record_id ) ;
            byte[] bytes = found.getBytes() ;
            String jsonText = new String(bytes) ;
            System.out.println( "Document Found: " + key ) ;    
            return jsonText ;
        } catch (SM.NotFoundException nfe) {
        	System.out.println( "Document Found: " + key ) ;    
			throw new DocumentException( "Document Not Found: " + key ) ;   
		} catch (Exception e) {
			throw new DocumentException( e.toString() ) ;                 
        }   	
    }


    public static SyncRequest get_sync_request(String key) throws DocumentException {
        System.out.println( "Get Document: " + key ) ;
        Document doc = KEYMAP_CACHE.get( key ) ;
        if ( doc == null || doc.record == null )
        {
            doc = DELETEKEY_CACHE.get( key ) ;
            if( doc == null)
            {
                throw new DocumentException( "Document Not Found: " + key ) ;
            }
            else
            {
                return get_sync_request(key, "delete") ;
            }
        }
        String record_key = doc.record ;
        SM db = SMFactory.getInstance() ;
        SM.OID record_id ;
        SM.Record found ;
        record_id = db.getOID( record_key.getBytes() ) ;
        try {
            found = db.fetch( record_id ) ;
            byte[] bytes = found.getBytes() ;
            String jsonText = new String(bytes) ;
            System.out.println( "Document Found: " + key ) ;    
            SyncRequest sync = new SyncRequest() ;
            sync.key = doc.key ;
            sync.json = jsonText ;
            sync.vclock[0] = (doc.vclock[0] == null) ? null : new String(doc.vclock[0]) ;
            sync.vclock[1] = (doc.vclock[1] == null) ? null : new String(doc.vclock[1]) ;
            sync.vclock[2] = (doc.vclock[2] == null) ? null : new String(doc.vclock[2]) ;
            sync.vclock[3] = (doc.vclock[3] == null) ? null : new String(doc.vclock[3]) ;
            sync.vclock[4] = (doc.vclock[4] == null) ? null : new String(doc.vclock[4]) ;
            sync.vclock[5] = (doc.vclock[5] == null) ? null : new String(doc.vclock[5]) ;
            sync.command = "" ; // set by caller
            return sync ;
        } catch (SM.NotFoundException nfe) {
            System.out.println( "Document Found: " + key ) ;    
            throw new DocumentException( "Document Not Found: " + key ) ;   
        } catch (Exception e) {
            throw new DocumentException( e.toString() ) ;                 
        }       
    }

    public static SyncRequest get_sync_request(String key, String command) throws DocumentException {
        System.out.println( "Get Document: " + key ) ;
        Document doc = DELETEKEY_CACHE.get( key ) ;
        if ( doc == null )
            throw new DocumentException( "Document Not Found: " + key ) ;
        try {
                SyncRequest sync = new SyncRequest() ;
                sync.key = doc.key ;
                sync.vclock[0] = (doc.vclock[0] == null) ? null : new String(doc.vclock[0]) ;
                sync.vclock[1] = (doc.vclock[1] == null) ? null : new String(doc.vclock[1]) ;
                sync.vclock[2] = (doc.vclock[2] == null) ? null : new String(doc.vclock[2]) ;
                sync.vclock[3] = (doc.vclock[3] == null) ? null : new String(doc.vclock[3]) ;
                sync.vclock[4] = (doc.vclock[4] == null) ? null : new String(doc.vclock[4]) ;
                sync.vclock[5] = (doc.vclock[5] == null) ? null : new String(doc.vclock[5]) ;
                sync.command = "" ; // set by caller
                return sync ;
            } catch (Exception e) {
                throw new DocumentException( e.toString() ) ;                 
            }       
        }


    public static void update_document( String key, String value ) throws DocumentException {
    	System.out.println( "update_document() Get Document: " + key + ", value: " + value ) ;
    	Document doc = KEYMAP_CACHE.get( key ) ;
    	if ( doc == null || doc.record == null )
    		throw new DocumentException( "Document Not Found: " + key ) ;
    	String record_key = doc.record ;
    	SM db = SMFactory.getInstance() ;
        SM.Record found ;
        SM.Record record ;
        SM.OID update_id ;        
		SM.OID record_id = db.getOID( record_key.getBytes() ) ;
		String jsonText = value ;
		int size = jsonText.getBytes().length ;
 		try {
            // store json to db
            record = new SM.Record( size ) ;
            record.setBytes( jsonText.getBytes() ) ;
            update_id = db.update( record_id, record ) ;
            System.out.println( "Document Updated: " + key ) ;
            // update vclock
            AdminServer server = AdminServer.getInstance() ;
            String my_host = server.getMyHostname() ;
            doc.vclock[0] = my_host ;
            int my_index = server.nodeIndex( my_host ) ;
            String old_version = doc.vclock[my_index] ;
            String[] splits = old_version.split(":") ;
            int version = Integer.parseInt(splits[1])+1 ;
            String new_version = my_host + ":" + Integer.toString(version) ;            
            doc.vclock[my_index] = new_version ;
            // docUpdaterNode.put(key,my_index) ;
            // sync nodes
            AdminServer.syncDocument( key, "update" ) ; 
			return ;             
        } catch (SM.NotFoundException nfe) {
			throw new DocumentException( "Document Not Found: " + key ) ;
       	} catch (Exception e) {
           	throw new DocumentException( e.toString() ) ;           
        }
    }


    public static void delete_document( String key ) throws DocumentException {
    	System.out.println( "Delete Document: " + key ) ;
    	Document doc = KEYMAP_CACHE.get( key ) ;
    	if ( doc == null || doc.record == null )
    		throw new DocumentException( "Document Not Found: " + key ) ;
        Document delDoc = new Document() ;
        delDoc.key = (doc.key == null ) ? null : new String(doc.key) ;
        delDoc.vclock[0] = (doc.vclock[0] == null) ? null : new String(doc.vclock[0]) ;
        delDoc.vclock[1] = (doc.vclock[1] == null) ? null : new String(doc.vclock[1]) ;
        delDoc.vclock[2] = (doc.vclock[2] == null) ? null : new String(doc.vclock[2]) ;
        delDoc.vclock[3] = (doc.vclock[3] == null) ? null : new String(doc.vclock[3]) ;
        delDoc.vclock[4] = (doc.vclock[4] == null) ? null : new String(doc.vclock[4]) ;
        delDoc.vclock[5] = (doc.vclock[5] == null) ? null : new String(doc.vclock[5]) ;
        delDoc.json = null ;
    	String record_key = doc.record ;
    	SM db = SMFactory.getInstance() ;
        SM.Record found ;
        SM.Record record ;     
		SM.OID record_id = db.getOID( record_key.getBytes() ) ;
       	try {
            AdminServer server = AdminServer.getInstance() ;
            String my_host = server.getMyHostname() ;
            delDoc.vclock[0] = my_host ;
            int my_index = server.nodeIndex( my_host ) ;
            String old_version = delDoc.vclock[my_index] ;
            String[] splits = old_version.split(":") ;
            int version = Integer.parseInt(splits[1])+1 ;
            String new_version = my_host + ":" + Integer.toString(version) ;            
            delDoc.vclock[my_index] = new_version ;
            DELETEKEY_CACHE.put( key, delDoc) ;
            db.delete( record_id ) ;
            // remove key map
            KEYMAP_CACHE.remove( key ) ;
            // sync nodes
            AdminServer.syncDocument( key, "delete" ) ; 
			System.out.println( "Document Deleted: " + key ) ;
        } catch (SM.NotFoundException nfe) {
           throw new DocumentException( "Document Not Found: " + key ) ;
        } catch (Exception e) {
         	throw new DocumentException( e.toString() ) ;            
        }		
    }

}




