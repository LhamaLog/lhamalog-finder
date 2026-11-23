import java.io.*;
import java.nio.charset.StandardCharsets;
import java.util.*;
import java.util.regex.*;

// Motor de referência Java, sem dependências. Protocolo stdin: regexes e
// marcadores em Base64, seguidos de uma linha Base64 por caso de conformidade.
// A Sugestão deverá usar o mesmo artefato e satisfazer estes vetores.
class Conformidade {
 static String ler(String s){return new String(Base64.getDecoder().decode(s),StandardCharsets.UTF_8);}
 static List<int[]> matches(Pattern p,String s){List<int[]> r=new ArrayList<>();Matcher m=p.matcher(s);while(m.find())r.add(new int[]{m.start(),m.end()});return r;}
 public static void main(String[] args)throws Exception{
  BufferedReader in=new BufferedReader(new InputStreamReader(System.in,StandardCharsets.UTF_8));
  Pattern email=Pattern.compile(ler(in.readLine()),Pattern.UNICODE_CASE);
  Pattern segredo=Pattern.compile(ler(in.readLine()),Pattern.UNICODE_CASE);
  String placeholder=ler(in.readLine()),marcador=ler(in.readLine());String line;
  while((line=in.readLine())!=null){
   String s=ler(line);List<int[]> found=matches(email,s);StringBuilder out=new StringBuilder();int pos=0;
   for(int[] m:found){out.append(s,pos,m[0]).append(placeholder);pos=m[1];}out.append(s,pos,s.length());
   int secrets=matches(segredo,out.toString()).size();boolean removida=secrets>0;
   String texto=removida?marcador:out.toString();
   System.out.println(Base64.getEncoder().encodeToString(texto.getBytes(StandardCharsets.UTF_8))+"\t"+(found.size()+secrets)+"\t"+removida);
  }
 }
}
