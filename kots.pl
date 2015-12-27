#!/usr/bin/env perl

use DateTime qw();
use utf8;

use strict;
use warnings;

# midnight = PST 7am
my $yesterday = DateTime->now->subtract(days => 1);
my $date = $yesterday->strftime("%y%m%d");
my $host = '192.168.2.100:9091';

run_tranmission('무한도전');
run_tranmission('런닝맨');
run_tranmission('냉장고를');

sub run_tranmission {
  my ($show) = @_;
  my $magnet_cmd = "go run kots.go -regex '^$show.+$date\.HDTV\.H264\.720p-WITH\$' -show $show";
  print "running $magnet_cmd\n";
  my $magnet_link =`$magnet_cmd`;
  return unless $magnet_link;
  print "found $magnet_link\n";

  print `transmission-remote "$host" -n transmission:transmission -a "$magnet_link"`, "\n";
}
