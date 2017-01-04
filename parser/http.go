package parser

import "gopkg.in/mgo.v2/bson"
import "net/url"
import "strings"

// HTTP provides a data structure for entries in bro's HTTP log file
type HTTP struct {
	// ID is the object id as set by mongodb
	ID bson.ObjectId `bson:"_id,omitempty"`
	// TimeStamp of this connection
	TimeStamp int64 `bson:"ts" bro:"ts" brotype:"time"`
	// Uid is the Unique Id for this connection (generated by Bro)
	UID string `bson:"uid" bro:"uid" brotype:"string"`
	// Source is the source address for this connection
	Source string `bson:"id_origin_h" bro:"id.orig_h" brotype:"addr"`
	// SourcePort is the source port of this connection
	SourcePort int `bson:"id_origin_p" bro:"id.orig_p" brotype:"port"`
	// Destination is the destination of the connection
	Destination string `bson:"id_resp_h" bro:"id.resp_h" brotype:"addr"`
	// DestinationPort is the port at the destination host
	DestinationPort int `bson:"id_resp_p" bro:"id.resp_p" brotype:"port"`
	// Transdepth is the ordinal value of requests into a pipeline transaction
	TransDepth int64 `bson:"trans_depth" bro:"trans_depth" brotype:"count"`
	// Version is the value of version in the request
	Version string `bson:"version" bro:"version" brotype:"string"`
	// Method is the request method used
	Method string `bson:"method" bro:"method" brotype:"string"`
	// Host is the value of the HOST header
	Host string `bson:"host" bro:"host" brotype:"string"`
	// URI is the uri used in this request
	URI string `bson:"uri" bro:"uri" brotype:"string"`
	// Referrer is the value of the referrer header in the request
	Referrer string `bson:"referrer" bro:"referrer" brotype:"string"`
	// UserAgent gives the user agent from the request
	UserAgent string `bson:"user_agent" bro:"user_agent" brotype:"string"`
	// ReqLen holds the length of the request body uncompressed
	ReqLen int64 `bson:"request_body_len" bro:"request_body_len" brotype:"count"`
	// RespLen hodls the length of the response body uncompressed
	RespLen int64 `bson:"response_body_len" bro:"response_body_len" brotype:"count"`
	// StatusCode holds the status result
	StatusCode int64 `bson:"status_code" bro:"status_code" brotype:"count"`
	// StatusMsg contains a string status message returned by the server
	StatusMsg string `bson:"status_msg" bro:"status_msg" brotype:"string"`
	// InfoCode holds the last seen 1xx informational reply code
	InfoCode int64 `bson:"info_code" bro:"info_code" brotype:"count"`
	// InfoMsg holds the last seen 1xx message string
	InfoMsg string `bson:"info_msg" bro:"info_msg" brotype:"string"`
	// FileName contains the name of the requested file
	FileName string `bson:"filename" bro:"filename" brotype:"string"`
	// Tags contains a set of indicators of various attributes related to a particular req and
	// response pair
	Tags string `bson:"tags" bro:"tags" brotype:"set[enum]"`
	// UserName will contain a username in the case of basic auth implementation
	UserName string `bson:"username" bro:"username" brotype:"string"`
	// Password will contain a password in the case of basic auth implementation
	Password string `bson:"password" bro:"password" brotype:"string"`
	// Proxied contains all headers that indicate a request was proxied
	Proxied string `bson:"proxied" bro:"proxied" brotype:"set[string]"`
	// OrigFuids contains an ordered vector of uniq file IDs
	OrigFuids string `bson:"orig_fuids" bro:"orig_fuids" brotype:"vector[string]"`
	// OrigFilenames contains an ordered vector of filenames from the client
	OrigFilenames string `bson:"orig_filenames" bro:"orig_filenames" brotype:"vector[string]"`
	// OrigMimeTypes contains an ordered vector of mimetypes
	OrigMimeTypes string `bson:"orig_mime_types" bro:"orig_mime_types" brotype:"vector[string]"`
	// RespFuids contains an ordered vector of unique file IDs in the response
	RespFuids string `bson:"resp_fuids" bro:"resp_fuids" brotype:"vector[string]"`
	// RespFilenames contains an ordered vector of unique files in the response
	RespFilenames string `bson:"resp_filenames" bro:"resp_filenames" brotype:"vector[string]"`
	// RespMimeTypes contains an ordered vector of unique MIME entities in the HTTP response body
	RespMimeTypes string `bson:"resp_mime_types" bro:"resp_mime_types" brotype:"vector[string]"`
}

// processHTTP fixes up absolute uri's as read by bro to be relative
func processHTTP(in interface{}) {
	line, found := in.(*HTTP)
	if !found {
		//this is the equivalent to a compile error
		panic("An object that is not *HTTP was passed into processHTTP")
	}

	//uri is missing the protocol. set uri to ""
	// ex: Host: 67.217.65.244 URI: 67.217.65.244:443
	if strings.HasPrefix(line.URI, line.Host) {
		line.URI = ""
		return
	}
	parsedURL, err2 := url.Parse(line.URI)
	if err2 != nil {
		line.URI = ""
		return
	}
	if parsedURL.IsAbs() {
		line.URI = parsedURL.RequestURI()
	}
}

// GetHostName is our method for collecting host name
func (in *HTTP) IsWhiteListed(whitelist []string) bool {
	if whitelist == nil {
		return false
	}
	if in.Host == "" {
		return false
	}

	for count := range whitelist {
		if strings.Contains(in.Host, whitelist[count]) {
			return true
		}
	}
	return false
}
