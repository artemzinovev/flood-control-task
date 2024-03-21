Когда завершите задачу, в этом README опишите свой ход мыслей: как вы пришли к решению, какие были варианты и почему выбрали именно этот. 

# Что нужно сделать

Реализовать интерфейс с методом для проверки правил флуд-контроля. Если за последние N секунд вызовов метода Check будет больше K, значит, проверка на флуд-контроль не пройдена.

- Интерфейс FloodControl располагается в файле main.go.

- Флуд-контроль может быть запущен на нескольких экземплярах приложения одновременно, поэтому нужно предусмотреть общее хранилище данных. Допустимо использовать любое на ваше усмотрение. 

# Необязательно, но было бы круто

Хорошо, если добавите поддержку конфигурации итоговой реализации. Параметры — на ваше усмотрение.

# Ход решения
Итак мной был реализован алгоритм по флуд-контролю. В первую очередь я обдумал общую модель логики алгоритма, каким образом можно определить временной интервал обращения к методу check и реализовать счетчик обращений. Реализовал алгоритм в рамках только языка го со сканером и мапой для определения работоспособности алгоритма.
Далее мной было выбрано хранилище данных redis, так как для реализации алгоритма основанном на хранении данных в качестве пары ключ-значение redis - это оптимальное решение. Далее в check было реализовано взаимодействие хранилища данных и алгоритма флуд-контроля.
После была реализована поддержка конфигурации итоговой реализации через конструктор CheckFloodConstructor.