package handler

import (
	"context"
	"log"

	"github.com/RomanGhost/buratino_bot.git/internal/handler/bot/function"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func InfoAboutInline(ctx context.Context, b *bot.Bot, update *models.Update) {
	function.InlineAnswer(ctx, b, update)

	inlineKeyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Создать ключ", CallbackData: ExtendKey},
			},
		},
	}

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

	inlineKeyboard := &models.InlineKeyboardMarkup{
		InlineKeyboard: [][]models.InlineKeyboardButton{
			{
				{Text: "Создать ключ", CallbackData: ExtendKey},
			},
		},
	}

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
