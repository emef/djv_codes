SRC_PATH="github.com/emef/djv_codes"

echo "deploying to $DJV_HOST"

echo "stopping server for deploy"
ssh $DJV_HOST "sudo systemctl stop djv_codes"

echo "uploding and installing djv_code_server"
ssh $DJV_HOST "mkdir -p go/src/github.com/emef/djv_codes"
rsync -az $GOPATH/src/github.com/emef/djv_codes $DJV_HOST:go/src/github.com/emef/
ssh $DJV_HOST "GOPATH=\$HOME/go go install github.com/emef/djv_codes/djv_code_server"
ssh $DJV_HOST "sudo cp go/bin/djv_code_server /usr/local/bin"
ssh $DJV_HOST "sudo chmod a+x /usr/local/bin/djv_code_server"
ssh $DJV_HOST "sudo mkdir -p /opt/djv_codes/codes"
ssh $DJV_HOST "sudo touch /opt/djv_codes/used_codes.txt"
ssh $DJV_HOST "sudo chown -R djv_codes /opt/djv_codes"

echo "installing systemd config and starting server"
scp $GOPATH/src/github.com/emef/djv_codes/deploy/djv_codes.service $DJV_HOST:
ssh $DJV_HOST "sudo mv djv_codes.service /lib/systemd/system/"
ssh $DJV_HOST "sudo chmod 755 /lib/systemd/system/djv_codes.service"
ssh $DJV_HOST "sudo systemctl daemon-reload"
ssh $DJV_HOST "sudo systemctl enable djv_codes.service"
ssh $DJV_HOST "sudo systemctl start djv_codes"

echo "installing nginx config and restarting it"
scp $GOPATH/src/github.com/emef/djv_codes/deploy/nginx.conf $DJV_HOST:
ssh $DJV_HOST "sudo rm /etc/nginx/sites-enabled/*"
ssh $DJV_HOST "sudo cp nginx.conf /etc/nginx/sites-available/djv_codes"
ssh $DJV_HOST "sudo ln -s /etc/nginx/sites-available/djv_codes /etc/nginx/sites-enabled/"
ssh $DJV_HOST "sudo service nginx restart"
