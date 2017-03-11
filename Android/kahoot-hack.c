#include<stdio.h>
#include<string.h>
int main(){
int hknm,pass;
char Script[20],Count[50],Pin[10],Space[2]=" ",Command[200],NameList[100],Name[20],UseNameList[1],yes[2]="y",Yes[2]="Y";
printf("\n1. kahoot-flood - using an old school denial of service technique, this program automatically joins a game of kahoot an arbitrary number of times. For instance, you can register the nicknames \"alex1\", \"alex2\", ..., \"alex100\".\n\n");
printf("2. kahoot-rand - connect to a game an arbitrary number of times (e.g. 100) and answer each question randomly. If you connect with enough names, one of them is bound to win.\n\n");
printf("3. kahoot-profane - circumvent Kahoot's profanity detector, allowing you to join with any nickname (but with extra length restrictions; it has to be short).\n\n");
printf("4. kahoot-play - play kahoot regularly-as if you were using the online client.\n\n");
printf("5. kahoot-html - I have notified Kahoot and they have fixed this issue. It used to allow you to join a game of kahoot a bunch of times with HTML-rich nicknames. This messes with the lobby of a kahoot game.\n\n");
printf("6. kahoot-crash - trigger an exception on the host's computer. This no longer prevents the game from functioning, so it is a rather pointless \"hack\"\n\n");
printf("7. kahoot-xss - since I discovered this security hole, I contacted Kahoot and they fixed it. This used to run arbitrary JavaScript code on the host's computer. This exploited a bug with the pre-game player list, which did not sanitize HTML tags. The exploit itself was rather complicated due to the fact that nicknames are limited to 15 characters.\n\n");
do{
  printf("Which hack do you want to use? [1-7]: ");
  scanf("%i",&hknm);
  if (hknm>0&&hknm<8){
    pass=1;
  }
}while (pass!=1);
if (hknm==1){
  printf("Game pin [123456]: ");
  scanf("%s",&Pin);
  printf("Use custom name list? [y/N]: ");
  scanf("%s",&UseNameList);
  if (strcmp(UseNameList,yes)==0||strcmp(UseNameList,Yes)==0){
    printf("Filename [name-list.txt]: ");
    scanf("%s",&NameList);
    strcpy(Command,"go run /data/data/com.termux/files/usr/var/kahoot-hack/flood.go ");
    strcat(Command,Pin);
    strcat(Command,Space);
    strcat(Command,NameList);
    system(Command);
  }else{
    printf("Nickname prefix [John]: ");
    scanf("%s",&Name);
    printf("Count [20]: ");
    scanf("%s",&Count);
    strcpy(Command,"go run /data/data/com.termux/files/usr/var/kahoot-hack/flood.go ");
    strcat(Command,Pin);
    strcat(Command,Space);
    strcat(Command,Name);
    strcat(Command,Space);
    strcat(Command,Count);
    system(Command);
  }
  }else if (hknm==2){
    printf("Game pin [123456]: ");
    scanf("%s",&Pin);
    printf("Use custom name list? [y/N]: ");
    scanf("%s",&UseNameList);
    if (strcmp(UseNameList,yes)==0||strcmp(UseNameList,Yes)==0){
      printf("Filename [name-list.txt]: ");
      scanf("%s",&NameList);
      strcpy(Command,"go run /data/data/com.termux/files/usr/var/kahoot-hack/rand.go ");
      strcat(Command,Pin);
      strcat(Command,Space);
      strcat(Command,NameList);
      system(Command);
    }else{
      printf("Nickname prefix [John]: ");
      scanf("%s",&Name);
      printf("Count [20]: ");
      scanf("%s",&Count);
      strcpy(Command,"go run /data/data/com.termux/files/usr/var/kahoot-hack/rand.go ");
      strcat(Command,Pin);
      strcat(Command,Space);
      strcat(Command,Name);
      strcat(Command,Space);
      strcat(Command,Count);
      system(Command);
    }
  }else if (hknm==3){
    printf("Game pin [123456]: ");
    scanf("%s",&Pin);
    printf("Nickname prefix [John]: ");
    scanf("%s",&Name);
    printf("Count: ");
    scanf("%s",&Count);
    strcpy(Command,"go run /data/data/com.termux/files/usr/var/kahoot-hack/profane.go ");
    strcat(Command,Pin);
    strcat(Command,Space);
    strcat(Command,Name);
    strcat(Command,Space);
    strcat(Command,Count);
    system(Command);
  }else if (hknm==4){
    printf("Game pin [123456]: ");
    scanf("%s",&Pin);
    printf("Nickname [John]: ");
    scanf("%s",&Name);
    strcpy(Command,"go run /data/data/com.termux/files/usr/var/kahoot-hack/play.go ");
    strcat(Command,Pin);
    strcat(Command,Space);
    strcat(Command,Name);
    system(Command);
  }else if (hknm==5){
    printf("Game pin [123456]: ");
    scanf("%s",&Pin);
    printf("Nickname [John]: ");
    scanf("%s",&Name);
    strcpy(Command,"go run /data/data/com.termux/files/usr/var/kahoot-hack/html.go ");
    strcat(Command,Pin);
    strcat(Command,Space);
    strcat(Command,Name);
    system(Command);
  }else if (hknm==6){
    printf("Game pin [123456]: ");
    scanf("%s",&Pin);
    printf("Nickname [John]: ");
    scanf("%s",&Name);
    strcpy(Command,"go run /data/data/com.termux/files/usr/var/kahoot-hack/crash.go ");
    strcat(Command,Pin);
    strcat(Command,Space);
    strcat(Command,Name);
    system(Command);
  }else if (hknm==7){
    printf("Game pin [123456]: ");
    scanf("%s",&Pin);
    printf("Script: ");
    scanf("%s",&Script);
    strcpy(Command,"go run /data/data/com.termux/files/usr/var/kahoot-hack/xss.go ");
    strcat(Command,Pin);
    strcat(Command,Space);
    strcat(Command,Script);
    system(Command);
  }
}
