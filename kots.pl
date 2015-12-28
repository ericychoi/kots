#!/usr/bin/env perl

# add to cron
# cd $WORKDIR && ./kots.pl

use POSIX qw(strftime);
use utf8;

use strict;
use warnings;

# midnight = PST 7am
my $date = strftime('%y%m%d', localtime);
my $host = '192.168.2.100:9091';

run_tranmission('무한도전');
run_tranmission('런닝맨');
#TODO: this doesn't work because drama uses different link, kots.go needs to change
#run_tranmission('응답하라 1988', "go run kots.go -regex '\[tvN\] 응답하라 1988\.E\\d\+\.$date\.HDTV\.Film\.x264\.720p-AAA' -show '응답하라 1988'");
run_tranmission('냉장고를');

sub run_tranmission {
  my ($show, $cmd) = @_;
  my $magnet_cmd = $cmd || "go run kots.go -regex '$show.+$date\.HDTV\.H264\.720p-WITH' -show $show";
  print "running $magnet_cmd\n";
  my $magnet_link =`$magnet_cmd`;
  return unless $magnet_link;
  print "found $magnet_link\n";

  print `transmission-remote "$host" -n transmission:transmission -a "$magnet_link"`, "\n";
}
