local http = require("socket.http")
local json = require("json")
local redirList = {}
local lastRedirUpdate = 0

function prequery ( dnspacket )
   pdnslog ("prequery called for ".. tostring(dnspacket) )
   remote_ip = dnspacket:getRemote()
   pdnslog ("prequery ip ".. remote_ip )
   qname, qtype = dnspacket:getQuestion()
   pdnslog ("q: ".. qname.." "..qtype)
   if lastRedirUpdate + 60 < os.time() then
      body, statuscode, headers = http.request("http://127.0.0.1:63636/redir/list")
      if pcall(function() getIpList(body) end) then
         pdnslog ("updating IP list")
         lastRedirUpdate = os.time()
      else
         pdnslog ("failed to update IP list")
      end
   end

   if qtype == pdns.A and redirList[remote_ip]
   then
      pdnslog ("calling dnspacket:setRcode")
      --      dnspacket:setRcode(pdns.NXDOMAIN)
      pdnslog ("called dnspacket:setRcode")
      pdnslog ("adding records")
      ret = {}
      ret[1] = {qname=qname, qtype=qtype, content=redirList[remote_ip], place=2, ttl=10}
      ret[2] = {qname=qname, qtype=pdns.TXT, content=os.date("Retrieved at %Y-%m-%d %H:%M:%S"), ttl=10}
      dnspacket:addRecords(ret)
      pdnslog ("returning true")
      return true
   end
   pdnslog ("returning false")
   return false
end


function getIpList (content)
   d = json.decode(content)
   if d then
      redirList = d
   end
end
