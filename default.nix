with import <nixpkgs> {};
#{ stdenv, buildGoModule, fetchFromGitHub }:
buildGoModule rec {
  name = "mqtt2prom";
  version = "1.0.1";
  src = fetchFromGitHub {
    owner = "jiphex";
    repo = name;
    rev = "v${version}";
    sha256 = "15937y6iwlwf0c4d1sikxikfm8jg6bgp27nfawxm1p987z6d0lfv";
  };
  vendorSha256 = "15ixxfx7cvb7xclwckjd3lif6qgjvqpzrfd853rs298i6mnxp4qm";
  subPackages = [ "cmd/mqtt2prom" ];
  runVend = false;
  meta = with stdenv.lib; {
    
  };
}
# stdenv.mkDerivation rec {
#   pname = "mqtt2prom";
#   version = "1.0.1";
#   
#   buildInputs = [ git go stdenv glibc.static ];
#
#   src = fetchFromGitHub {
#    owner = "jiphex";
#    repo = pname;
#    rev = "v${version}";
#    sha256 = "0m2fzpqxk7hrbxsgqplkg7h2p7gv6s1miymv3gvw0cz039skag0s";
#  };
#
#  runVend = false;
#
#  subpackages = [ "cmd/mqtt2prom" ];
#
#  meta = with lib; {
#    description = "MQTT Exporter";
#    homepage = "https://gitHub.com/jiphex/${pname}";
#    license = licenses.mit;
#    maintainers = with maintainers; [ ];
#    platforms = platforms.linux ;
#  };
#}
