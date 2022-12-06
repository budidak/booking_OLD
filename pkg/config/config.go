package config

import (
	"html/template"
	"log"

	"github.com/alexedwards/scs/v2"
)

type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
}

// Bu dosyada sadece programımızla ilgili konfigürasyon değişkenleri olacaktır. Bu yüzden bu dosyaya kendi yazdığımız programdan herhangi bir import yapmamalıyız (import cycle oluşursa sürekli paketler birbirini import eder ve programımız compile etmez, burada sadece gerekli olduğunda standart libraryde bulunan built-in packages import edilmelidir)

// Böyle bir şey yazmamızın nedeni, main() içerisinden bu ayarlarla oynayarak programımız üzerinde daha fazla kontrole sahip olmak. Örneğin development mode'da UseCache false haline getirip her şeyi diskten okuyabiliriz. Ama production mode'a geçtiğimizde bunu true yaparız ve eğer sayfa cache değişkenimizde kayıtlıysa bilgiyi oradan çekeriz.

// Buna REPOSITORY pattern ismi verilir.
