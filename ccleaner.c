#include <stdio.h>
#include <unistd.h>
#include <stdlib.h>
#include <ctype.h>
#include <stddef.h>
#include <stdbool.h>
#include <string.h>

bool validRecord(char *buf, int count, int lookahead_count) {
	if(count < 20) {
		return false;
	}

	if(buf[0] == '\0'){
		return false;
	}

	for(int i = 20; i < 20+lookahead_count; i++){
		if(buf[i] != ' ') {
			return true;
		}
	}

	return false;
}

int main(int argc, char *argv[]) {

  char *infile = NULL;
  char *outfile = "out.dat";
  int record_length = 1024;
  int lookahead_count;
	int c;

  opterr = 0;

	while((c = getopt (argc, argv, "i:o:c:l:")) != -1)
    switch (c)
      {
      case 'i':
        infile = optarg;
        break;
      case 'o':
        outfile = optarg;
        break;
      case 'c':
        record_length = atoi(optarg);
        break;
      case 'l':
        lookahead_count = atoi(optarg);
        break;
      default:
        abort();
      }

	if (lookahead_count == 0) {
		lookahead_count = record_length - 20;
	}

	if(infile == NULL || outfile == NULL || record_length <= 20 || lookahead_count > (record_length-20)) {
		printf("Usage: -i <infile> -o <outfile> -c <record length, must be > 20> -l <lookahead, must be <= record length - 20");
		return 1;
	}

  printf("Copying from %s to %s with record length %d and lookahead %d\n",
          infile, outfile, record_length, lookahead_count);

	FILE *in;
	FILE *out;
	char buf[record_length];
	int numRecords = 0;
	int nullRecords = 0;

	in = fopen(infile, "r");
	out = fopen(outfile, "w");

	if(in == NULL) {
		printf("Could not open file %s for reading.\n", infile);
		return 1;
	}

	if(out == NULL) {
		printf("Could not open file %s for writing.\n", outfile);
		return 1;
	}

	int count = fread(buf, 1, record_length, in);

	fwrite(buf, 1, count, out);

	while(count > 0) {
		count = fread(buf, 1, record_length, in);

		if(count != 0 && validRecord(buf, count, lookahead_count)) {
			numRecords++;
			fwrite(buf, 1, count, out);
		} else {
			nullRecords++;
		}
	}

	printf("Copied %d records and skipped %d null records.", numRecords, nullRecords-1);
	fclose(in);
	fclose(out);

  return 0;
}

