package handler

import (
	"context"
	"log"

	"github.com/RomanGhost/buratino_bot.git/internal/telegram/data"
	"github.com/RomanGhost/buratino_bot.git/internal/telegram/function"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func InfoAboutInline(ctx context.Context, b *bot.Bot, update *models.Update) {
	function.InlineAnswer(ctx, b, update)

	inlineKeyboard := data.CreateKeyboard([]models.InlineKeyboardButton{data.CreateKeyButton()})

	message := `📜 *Сказ о волшебных ключах* 🗝️
В этой сказочной обители ты встретил *Буратино* \- не просто деревянного мальчишку, а стража потайных троп интернета\! 🌐✨
Он дарует *волшебные VPN\-ключи*, что действуют недолго \- всего около *30 минут*, но дают силу обойти коварных Карабасов и Брандмейстеров\.
Сейчас ключики выдаются через *Outline*, но скоро и *WireGuard* придёт на помощь храбрым странникам\!
Нажми на волшебную кнопку, и путь будет открыт\.\.\. 🧙‍♂️🔑`
	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        message,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: inlineKeyboard,
	})

	if err != nil {
		log.Printf("[WARN] Error send info message %v", err)
	}
}

func HelpOutlineIntructionInline(ctx context.Context, b *bot.Bot, update *models.Update) {
	function.InlineAnswer(ctx, b, update)

	inlineKeyboard := data.CreateKeyboard([]models.InlineKeyboardButton{data.CreateKeyButton()})

	message := `📜 *Волшебная инструкция по настройке VPN*
Следуй за мной, деревянный друг, в страну свободного интернета\! 🌍✨ Вот как обрести силу волшебного ключа:
🔧 *Шаг 1: Установи волшебное зеркало \- Outline App:*
📱 iOS: https://itunes\.apple\.com/app/outline\-app/id1356177741
🍏 MacOS: https://itunes\.apple\.com/app/outline\-app/id1356178125
🪟 Windows: https://s3\.amazonaws\.com/outline\-releases/client/windows/stable/Outline\-Client\.exe
🐧 Linux: https://s3\.amazonaws\.com/outline\-releases/client/linux/stable/Outline\-Client\.AppImage
🤖 Android: https://play\.google\.com/store/apps/details\?id\=org\.outline\.android\.client
🔄 Android \(альтернатива\): https://s3\.amazonaws\.com/outline\-releases/client/android/stable/Outline\-Client\.apk

🔑 *Шаг 2: Жди волшебный ключ\!* 
Ты получишь таинственный ключик, что начинается с \'ss://\' \- скопируй его, как только он появится\! ✨

🚪 *Шаг 3: Вставь ключ в Outline и открой врата свободы\!* 
Если приложение само распознает ключ \- просто нажми _Connect_\.
Если нет \- вставь его вручную и тоже нажми _Connect_\.

✅ *Готово\!* Чтобы убедиться, что ты в стране свободного интернета, загугли: _what is my ip_ и сравни IP с тем, что в Outline\.

🧙‍♂️ Пусть ни один Карабас не догонит тебя в этом цифровом приключении\!`

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        message,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: inlineKeyboard,
	})

	if err != nil {
		log.Printf("[WARN] Error send info message %v", err)
	}
}

func HelpWireguardIntructionInline(ctx context.Context, b *bot.Bot, update *models.Update) {
	function.InlineAnswer(ctx, b, update)

	inlineKeyboard := data.CreateKeyboard([]models.InlineKeyboardButton{data.CreateKeyButton()})

	message := `📜 *Волшебная инструкция по настройке VPN с WireGuard*

Следуй за мной, деревянный друг, в страну безопасного и свободного интернета\! 🌍✨ 
Вот как обрести силу волшебного туннеля:

🔧 *Шаг 1: Установи волшебное зеркало \- WireGuard App:*
📱 iOS: https://apps\.apple\.com/app/wireguard/id1441195209
🍏 MacOS: https://apps\.apple\.com/app/wireguard/id1451685025
🪟 Windows: https://www\.wireguard\.com/install/
🐧 Linux: https://www\.wireguard\.com/install/
🤖 Android: https://play\.google\.com/store/apps/details?id\=com\.wireguard\.android

🔑 *Шаг 2: Жди волшебный ключ\!* 
Ты получишь загадочный конфиг\-файл с расширением \'\.conf\' или текстовый ключ\. 
Сохрани его, как настоящий амулет ✨

🚪 *Шаг 3: Вставь ключ в WireGuard и открой врата свободы\!* 
\- Если у тебя файл \'\.conf\' — просто импортируй его в приложение\.  
\- Если текстовый ключ — создай новый туннель вручную, вставив публичный ключ сервера, 
  приватный ключ клиента, адреса и порты\.  
\- Нажми _Activate_ или _Connect_ и почувствуй магию соединения ✨

✅ *Готово\!* 
Чтобы убедиться, что ты в стране свободного интернета, загугли: _what is my ip_ 
и сравни IP с тем, что указан в твоем туннеле WireGuard\.

🧙‍♂️ Пусть ни один цифровой дракон не сможет преградить тебе путь\!`

	_, err := b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.CallbackQuery.Message.Message.Chat.ID,
		Text:        message,
		ParseMode:   models.ParseModeMarkdown,
		ReplyMarkup: inlineKeyboard,
	})

	if err != nil {
		log.Printf("[WARN] Error send info message %v", err)
	}
}
