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
my $host = '192.168.1.2:9091';

GetOptions( 'date=s' => \$date );

run_tranmission('무한도전');
run_tranmission('런닝맨');
run_tranmission('아는 형님');
run_tranmission('냉장고를');
#run_tranmission('비정상회담');
run_tranmission('너의 목소리가');
#run_tranmission('개그콘서트');
run_tranmission('마이 리틀 텔레비전');
run_tranmission('꽃보다 청춘');
run_tranmission('최고의 사랑');
run_tranmission('삼시세끼');
run_tranmission('SNL 코리아');
run_tranmission('신서유기');
run_tranmission('SHOW ME THE MONEY');
run_tranmission('노래의 탄생', "./kots -regex '노래의 탄생.+\\.$date\\.720p-NEXT' -show '노래의 탄생'");
#run_tranmission('진짜 사나이');
run_tranmission('정글의 법칙', "./kots -regex '정글의 법칙.+\\.$date\\..+720p-NEXT' -show '정글의 법칙'");
run_tranmission('해피 투게더');
#run_tranmission('언프리티 랩스타');
run_tranmission('슬램덩크');
run_tranmission('미운 우리 새끼');
run_tranmission('K팝스타');

sub run_tranmission {
  my ($show, $cmd) = @_;
  my $magnet_cmd = $cmd || "./kots -regex '$show.+\\.$date\\..*720p-NEXT' -show $show";
  my $encoded_cmd = encode("utf8", $magnet_cmd);
  print +localtime().": running $encoded_cmd\n";
  my $magnet_link = `$encoded_cmd`;
  return unless $magnet_link;
  print +localtime().": found $magnet_link, running transmission-remote\n";
  print +localtime().": ".`transmission-remote "$host" -n transmission:transmission -a "$magnet_link"`, "\n";
}
