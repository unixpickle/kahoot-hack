decode('hVKdIFfQ43UON4WvocX1onmx6cZSHsiG1BTdDahaepQej1nn6tu0m8B2b68zX74v2T5cV581Y7vNnU0RcIqTzHBOn9rPdqvKMv3z'); function decode(message) {var offset = (7 + (64 + 65) * (42 * 98) + 49); console.log("Offset derived as:", offset); return _.replace(message, /./g, function(char, position) {return String.fromCharCode((((char.charCodeAt(0) * position) + offset) % 77) + 48);});}

