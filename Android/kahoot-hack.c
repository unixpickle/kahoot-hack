#include<stdio.h>
#include<string.h>
int main(){
int hknm,pass;
char spt[20],nob[50],pin[10],spc[2]=" ",cmd[200],nls[100],nme[20],cnl[1],cnly[2]="y",cnlY[2]="Y";
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
  scanf("%s",&pin);
  printf("Use custom name list? [y/N]: ");
  scanf("%s",&cnl);
  if (strcmp(cnl,cnly)==0||strcmp(cnl,cnlY)==0){
    printf("Filename [name-list.txt]: ");
    scanf("%s",&nls);
    strcpy(cmd,"go run /data/data/com.termux/files/usr/var/kahoot-hack/flood.go ");
    strcat(cmd,pin);
    strcat(cmd,spc);
    strcat(cmd,nls);
    printf("System will exec this: %s",cmd);
  }else{
    printf("Nickname prefix [John]: ");
    scanf("%s",&nme);
    printf("Count [20]: ");
    scanf("%s",&nob);
    strcpy(cmd,"go run /data/data/com.termux/files/usr/var/kahoot-hack/flood.go ");
    strcat(cmd,pin);
    strcat(cmd,spc);
    strcat(cmd,nme);
    strcat(cmd,spc);
    strcat(cmd,nob);
    printf("System will exec this: %s",cmd);
  }
  }else if (hknm==2){
    printf("Game pin [123456]: ");
    scanf("%s",&pin);
    printf("Use custom name list? [y/N]: ");
    scanf("%s",&cnl);
    if (strcmp(cnl,cnly)==0||strcmp(cnl,cnlY)==0){
      printf("Filename [name-list.txt]: ");
      scanf("%s",&nls);
      strcpy(cmd,"go run /data/data/com.termux/files/usr/var/kahoot-hack/rand.go ");
      strcat(cmd,pin);
      strcat(cmd,spc);
      strcat(cmd,nls);
      printf("System will exec this: %s",cmd);
    }else{
      printf("Nickname prefix [John]: ");
      scanf("%s",&nme);
      printf("Count [20]: ");
      scanf("%s",&nob);
      strcpy(cmd,"go run /data/data/com.termux/files/usr/var/kahoot-hack/rand.go ");
      strcat(cmd,pin);
      strcat(cmd,spc);
      strcat(cmd,nme);
      strcat(cmd,spc);
      strcat(cmd,nob);
      printf("System will exec this: %s",cmd);
    }
  }else if (hknm==3){
    printf("Game pin [123456]: ");
    scanf("%s",&pin);
    printf("Nickname prefix [John]: ");
    scanf("%s",&nme);
    printf("Count: ");
    scanf("%s",&nob);
    strcpy(cmd,"go run /data/data/com.termux/files/usr/var/kahoot-hack/profane.go ");
    strcat(cmd,pin);
    strcat(cmd,spc);
    strcat(cmd,nme);
    strcat(cmd,spc);
    strcat(cmd,nob);
    printf("System will exec this: %s",cmd);
  }else if (hknm==4){
    printf("Game pin [123456]: ");
    scanf("%s",&pin);
    printf("Nickname [John]: ");
    scanf("%s",&nme);
    strcpy(cmd,"go run /data/data/com.termux/files/usr/var/kahoot-hack/play.go ");
    strcat(cmd,pin);
    strcat(cmd,spc);
    strcat(cmd,nme);
    printf("System will exec this: %s",cmd);
  }else if (hknm==5){
    printf("Game pin [123456]: ");
    scanf("%s",&pin);
    printf("Nickname [John]: ");
    scanf("%s",&nme);
    strcpy(cmd,"go run /data/data/com.termux/files/usr/var/kahoot-hack/html.go ");
    strcat(cmd,pin);
    strcat(cmd,spc);
    strcat(cmd,nme);
    printf("System will exec this: %s",cmd);
  }else if (hknm==6){
    printf("Game pin [123456]: ");
    scanf("%s",&pin);
    printf("Nickname [John]: ");
    scanf("%s",&nme);
    strcpy(cmd,"go run /data/data/com.termux/files/usr/var/kahoot-hack/crash.go ");
    strcat(cmd,pin);
    strcat(cmd,spc);
    strcat(cmd,nme);
    printf("System will exec this: %s",cmd);
  }else if (hknm==7){
    printf("Game pin [123456]: ");
    scanf("%s",&pin);
    printf("Script: ");
    scanf("%s",&spt);
    strcpy(cmd,"go run /data/data/com.termux/files/usr/var/kahoot-hack/xss.go ");
    strcat(cmd,pin);
    strcat(cmd,spc);
    strcat(cmd,spt);
    printf("System will exec this: %s",cmd);
  }
}
