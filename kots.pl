#!/usr/bin/env perl

# add to cron
# cd $WORKDIR && ./kots.pl

use POSIX qw(strftime);
use utf8;
use Encode qw(encode);

use strict;
use warnings;

use Getopt::Long;

# midnight = PST 7am
my $date = strftime('%y%m%d', localtime);
my $host = '192.168.2.100:9091';

GetOptions( 'date=s' => \$date );

run_tranmission('무한도전', "./kots -regex '무한도전\\..+\\.$date\\.HDTV\\.H264\\.720p-Venus' -show '무한도전'");
run_tranmission('런닝맨');
run_tranmission('냉장고를');
run_tranmission('비정상회담');
run_tranmission('너의 목소리가');
run_tranmission('개그콘서트');
run_tranmission('마이 리틀 텔레비전');
run_tranmission('문제적 남자');
run_tranmission('꽃보다 청춘');
run_tranmission('님과 함께');
run_tranmission('SNL코리아');
run_tranmission('신서유기');
run_tranmission('SHOW ME THE MONEY');
run_tranmission('노래의 탄생', "./kots -regex '노래의 탄생.+\\.$date\\.720p-NEXT' -show '노래의 탄생'");
run_tranmission('유희열의 스케치북', "./kots -regex '유희열의 스케치북.+\\.$date\\.720p-NEXT' -show '유희열의 스케치북'");
run_tranmission('예림이네 만물트럭');
run_tranmission('1박2일');
run_tranmission('진짜 사나이');
run_tranmission('정글의 법칙', "./kots -regex '정글의 법칙.+\\.$date\\.720p-NEXT' -show '정글의 법칙'");
run_tranmission('슈가맨', "./kots -regex '슈가맨.+\\.$date\\.720p-NEXT' -show '슈가맨'");

sub run_tranmission {
  my ($show, $cmd) = @_;
  my $magnet_cmd = $cmd || "./kots -regex '$show.+$date\.HDTV\.H264\.720p-NEXT' -show $show";
  my $encoded_cmd = encode("utf8", $magnet_cmd);
  print +localtime().": running $encoded_cmd\n";
  my $magnet_link = `$encoded_cmd`;
  return unless $magnet_link;
  print +localtime().": found $magnet_link, running transmission-remote\n";
  print +localtime().": ".`transmission-remote "$host" -n transmission:transmission -a "$magnet_link"`, "\n";
}
