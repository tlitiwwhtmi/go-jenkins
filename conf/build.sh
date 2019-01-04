workspace=${WORKSPACE}
PACKAGE_PATH=`find $workspace/ -regex ".*[0-9]+\.\(war\|zip\|tar\)" | egrep "$(date +%Y%m%d)|$(date +%Y%m%d-)" | grep -v ".repository"`
deploy_app=`echo ${PACKAGE_PATH} | tr '\ ' '\n' | awk -F/ '{print $NF}'`

if [ -z "${PACKAGE_PATH}" ];then
  exit 1
else
	 start=$(date +%s)
         for i in $PACKAGE_PATH;do
	    #curl -C - -sT $i ftp://deploy:jtwmydtsgx@10.32.135.110/$(date +%Y%m%d)/
            lftp ftp://***:***@ftpserver/$(date +%Y%m%d)/<<!
                mput $i
                close
                bye
!
           curl -s 'http://'${backIp}':'${backPort}'/savepac?hid='${backId}'&pac='$i
         done
  sleep 30
  curl -X POST http://messeraddress/jq/messer/api/addproject/ -d "appname=$(echo ${deploy_app} | tr '\ ' '\,')" -u ****:****
finished=$(date +%s)
time_take=`echo $(($finished - $start))`
 echo  "It takes -----"$time_take"-----Second to finish transfering  Packages"
fi

