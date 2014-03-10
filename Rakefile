def root_dir
  File.dirname(__FILE__)
end

task :build_dir do
  sh "mkdir -p build/src/github.com/farcaller"
  sh "ln -sf #{root_dir} build/src/github.com/farcaller/dockdns"
end

task :build_deps => :build_dir do
  sh "env GOPATH=\"#{root_dir}/build\" go get github.com/fsouza/go-dockerclient github.com/miekg/dns"
end

desc "Build dockdns binary"
task :build, [:goos, :goarch] => [:build_dir, :build_deps] do |t, args|
  args.with_defaults(goos: 'linux', goarch: 'arm')

  sh "env GOOS=#{args[:goos]} GOARCH=#{args[:goarch]} GOPATH=\"#{root_dir}/build\" go build -o build/dockdns github.com/farcaller/dockdns"
end

desc "Deploy dockdns to server, reload/restart systemd service"
task :deploy, [:remote] => [:build] do |t, args|
  remote = args[:remote]
  sh "scp build/dockdns support/dockdns.service #{remote}:/tmp && " +
     "ssh #{remote} sudo 'bash -c \"" +
       "chown 0:0 /tmp/dockdns /tmp/dockdns.service && " +
       "mv /tmp/dockdns /usr/bin/dockdns && " +
       "mv /tmp/dockdns.service /usr/lib/systemd/system/dockdns.service && " +
       "systemctl daemon-reload && " +
       "systemctl restart dockdns" +
       "\"'"
end

task default: :deploy
