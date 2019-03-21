# PDU 一种基于去中心化账户系统的社交网络
Parallel Digital Universe - A decentralized identity-based social network

email: hello@pdu.pub
微信: ![wechat](Wechat.png)

[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/pdupub/go-pdu)
[![GoReport](https://goreportcard.com/badge/github.com/pdupub/go-pdu)](https://goreportcard.com/report/github.com/pdupub/go-pdu)
[![Travis](https://travis-ci.org/pdupub/go-pdu.svg?branch=master)](https://travis-ci.org/pdupub/go-pdu)
[![License](https://img.shields.io/badge/license-GPL%20v3-blue.svg)](LICENSE)

**Abstract:** 通常意义上的社交网络服务(SNS)，如Facebook、twitter、微信等，用户可以在其上创建身份，维护关系并进行信息传播、交互。但现有的SNS均依赖于某个第三方提供的中心化的网络服务，这使得其很容易被控制或封锁隔离。BitTorrent协议，能够实现P2P的信息传播，但其根本目的是提高对于已知内容的传播效率，其弱化的账户系统设计不利于使用者对未知内容有所判别。对于去中心化的系统，即便用数字签名能够证明每个信息的来源，但是因为缺少第三方验证（如手机号注册）来给予账户创建成本，即便无用的信息充斥整个网络也无法信息来源进行惩罚。

我们提出一种在纯粹P2P的环境下给予账户创建成本的方式，并基于这种账户系统，构建完整的P2P社交网络形态。首先，我们引入的时间证明，用以证明某个特定行为发生于某时刻之后。然后，在系统中新账户的创建必须由多个已存在的合法账号联合签名，并规定同一账号的此类签名操作需满足时间间隔。每个账户系统的使用者（包括非用户），都在本地以DAG的结构维护部分或全部账号之间的关系拓扑，并随时可以根据自己获知的消息，对新的账号进行验证增补，同时也可因作恶行为对某些账号及关联账号进行惩罚。

与比特币为代表的区块链不同，系统的使用者并不需要对每个行为产生共识，维护同一个一致的信息。系统中的每个用户都只需要维护和自己使用相关的账户部分和信息，并且根据自己的标准对于关系拓扑中的账户做出是否接受其信息的判断。

<!-- MarkdownTOC depth=4 autolink=true bracket=round list_bullets="-*+" -->
- [Introduction](#introduction)
  * [现状](#现状) 
  * [PDU](#pdu)
- [Time Proof](#time-proof)
- [Account Topology](#account-topology) 
  * [亲源关系](#亲源关系)
  * [生命周期](#生命周期)
  * [自然法则](#自然法则)
- [Message](#message)
  * [Message Credit](#message-credit)
  * [Time Proof Message](#time-proof-message)
  * [Cosigned Birth Message](#cosigned-birth-message)
  * [Evidence Message](#evidence-message)
- [Network](#network)
  * [Message Spread](#message-spread)
  * [Account Create](#account-create)
  * [PDU Evolve Step](#pdu-evolve-step)
- [Function Node](#function-node)
  * [Time Proof Node](#time-proof-node)
  * [Account Node](#account-node)
  * [Tracker Node](#tracker-node)
  * [Message Node](#message-node)
- [Conclusion](#conclusion)
<!-- /MarkdownTOC -->

## Introduction

#### 现状

现今互联网上的信息传播、交互大多依赖于一个强大可信的第三方中心化服务，如Facebook、Twitter、微信、微博等社交网络服务。其存在毋庸置疑给使用者带来了极大的便利，但随着其逐步发展，中心化社交服务的问题也逐渐显露。

1. 无论有意或无意，第三方服务都存在越权使用用户信息或者造成用户数据泄露的可能。
2. 出于商业的考量，中心化的服务商会利用自身强大的用户基础，来打压竞品，如限制其产品的信息在自有平台上传播，以维护自身的垄断地位。
3. 中心化的服务容易受制于政府的监管，封锁。

但即便如此，由于对于三方中心化服务的依赖，很多用户依然不得已选择继续使用发生问题的服务，而非迁移自己的数据。因为对于大多数用户而言，离开某个平台虽然不会损失自己的数据信息，但却失去了在此平台上长期积累的用户关系和自身在此平台信用值。

本质上来说，这个问题根源在于用户群体自身并不能构成一个脱离某第三方的网络，所以用户的关系信息归属于平台而非其自身。

#### PDU

我们提出一种新的基于去中心化账户系统的社交网络（PDU）的本意也并排除第三方的中心化服务，而是希望能够通过去中心化账户系统（DID）的实现，能够将用户身份确认及关系拓扑脱离于某个特定平台，用以消除用户对于特定第三方中心化服务的依赖，让用户的身份及社交关系真正归属于用户。

如同*双花*可被认为是去中心化的货币系统需要解决的根本问题，一个去中心化的账户系统要解决的根本问题是如何给账户的创建赋予必要成本。

我们仿照自然及社会，首先引入时间证明的概念，并以此为基础订立自然法则。符合自然法则创建的新账户才有可能被系统中的其他用户（部分用户）所接受。每个用户自身都可用有向无环图（DAG）的结构来维护自身所*承认*的所有用户及其构成的亲源拓扑关系。任何违背自然法则的消息都会作为证据在网络中传播，让消息接受者可以根据本地的亲源拓扑对作恶的用户进行惩罚。惩罚的账户亲源深度、广度由接受者自行决定。

与传统中心化服务的账户系统不同，PDU的自然法则还基于时间证明定义了账户的生命周期，使得不被使用的账户可以被自然淘汰，账户的总数量呈线性增长（受时间流速的影响会），而当前生命周期内的用户数量会近似恒定。

时间证明是PDU中用户一切行为成本控制的基础，但因为PDU中没有强制的共识，取而代之的是用户自身的选择，所以完全可能有多个不同的时间证明的存在。PDU接受这种情况的存在，就如同平时存在的多个时空，甚至每个时空可以设定不同的时间流速来影响本时空中的行为成本。同时，任何用户也可同时存在于自己选择的多个时空当中。	

## Time Proof 

时间证明本身是一种消息（Message），并无特别。这种消息的每条信息中（*注意不是其时间证明字段中，而是信息内容中*）包含一个整数和一个字符串，数字作为这个时间证明当前的时间（*跟现实时间无关*），而字符串则作为这个时间证明的佐证（*通常为由时间戳+随机产生+消息内容之后计算hash，不可预测*）。这个结构连同这条消息的签名，会被一起放入采用这个时间证明的消息中，作为那条消息的产生时间晚于某个特定时刻的证明。

时间证明中除了包含上述基本内容之外也可以包含消息中的任意内容，如其他的时间证明（Time Proof），实时的观测数据，其他的消息内容等用以增强本时间证明的可信程度。需要注意的是不同的时间证明中所包含的时间可能被用于构建不同的时空，所以即便是统一消息包含的多个时间证明，其中的各个数值之间也可能毫无关系。只是用来让本消息在多个时空中均合法。同时，时间证明也未规定相邻两个块之间的时间间隔。

时间证明的选择权在每个用户，用户可以在行为（消息）中选择一个，多个，或者完全不选择任何时间证明，也可以在自己的多条消息中选择不同的时间证明。但推荐用户尽量选择可信度高的时间证明来为自己的行为设置时间证据，以防止由于发出时间证明的账号的作恶行为影响自己所采用的时间证明的公信力。因为时间证明本身属于消息，所以其可信度的判断规则同普通消息可信度的判断规则。

对于时间证明的发布者而言，创建新的时间证明可以从某个已存在的较为通用的时间证明的任意时间开始分叉，这样的做法可以尽可能保留更多的当前已经合法的用户。使得本证明更容易被更多的账号所采用。因为账号的创建成本，生命周期等都受制于时间，所以不同的时间证明可以对应不同的信息生命周期，很可以一个账号在某个时间线上生命周期已经结束，而在其他的时间线上还未结束。

## Account Topology

账户系统是用户在社交网络中一切行为的基础。基于账户，社交关系才得以建立，认证行为能够以发生，用户也才会因为自身的行为而得到奖惩的反馈。当账户系统由一个中心化的平台进行维护的时候，账号的创建过程，使用过程都基这个平台，所以很容易进行控制并有效的限制一些恶意的行为。比如为应单个使用者创建大量账号的行为，平台可通过手机号验证等绑定真实世界信息的方式来增加创建账号的成本。为应对身份冒用，盗取的行为，平台会在注册过程中强制要求用户使用更加复杂的密码，缩短登录的过期时间，加强自身平台的安全等级等方式。为应对用户的恶意行为，平台会指定一些规则条文，当用户触犯某些规则的时候，由平台对用户进行惩罚，这些惩罚的方式并不一定被用户所知晓，比如仅仅减少其信息的露出概率，又比如彻底删除其所有的信息。可见，对于账户系统的控制权利，完全在于其依赖的平台，当此平台完全可信的时候，这是一个很好的解决方案，但是否存在完全可信的中心化平台，答案是显然是否定的。

由于数字签名的存在，即便在一个完全P2P的网络环境中，对于信息的认证，保密等均不存在问题（*可优于中心化平台*）。用户的身份基于一个非对称的秘钥对，信息生产者利用私钥对信息进行签名，信息接受者用生产者公钥验证信息来源的真实性。对于加密内容，生产者可以接收者的公钥进行加密之后，再用生产者自身的私钥进行签名，信息接受者收到信息，先用对方公钥进行验证，再用自身的私钥对内容进行解密。

但由于作为身份基础的非对称秘钥对创建容易，单一使用者也可以在短时间内创建大量的秘钥对。为在P2P网络环境中，为了控制基于秘钥对的合法账户的数量，增加账户的创建成本，我们基于时间证明，首先提出亲源关系和生命周期两个概念，并在此基础上定义了多条自然法则。P2P网络中的每个用户，都可以依照其对于其他的用户进行判断，选择自己是否接受对方的存在。

#### 亲源关系

亲源关系是指两个账号之间的关系，在PDU中，每个节点所承认的单一时空的所有合法账户都存在直接或间接的亲源关系。除创世的两个账号之外，所有的账号均有且只有两个属性不同的父级账号。整个账户体系所构成的关系拓扑是由两个创世端点启始的有向无环图（DAG）

#### 生命周期

每个账户有其自己的以时间证明为基础的生命周期，这个账户的生命周期起始于两个异属性节点完成签名，并广播此事件的时间证明。一个账户只有在其生命周期之内生产的消息才能被认为合法。（*由于接受信息的节点会更倾向于时刻新且可信度高的消息，所以在生命周期结束以后，伪造历史消息进行广播的意义并不大。在Message章节中会具体叙述*）

生命周期的长度跟父级账号相关，但不低于某个特定值，那个特定值就是最低生命周期。


#### 自然法则

1. 每个账户都有一个二元属性，属性值以其公钥的末尾奇偶性确定，创世的两个账户必为异属性账户。
（*此规定意味着用户可以通过重复生成非对称秘钥对的过程来自己选择此二元化的属性，这个账户系统的次二元化属性不会趋近于统一，因为当一方变得稀有时，由于父级地址的签名规则，所以选择较稀有属性的地址会增多，以增加自身的账户价值。*）
2. 每个新创的非对称秘钥，需要进行签名过程，被两个合法异属性账户行签名，之后广播到整个P2P网络，才可能被其他账号所认可。
3. 签名执行的父级地址，在签名包含的时间戳之前，必须已经经历至少1/4个**最低生命周期**。
4. 签名执行的父级地址，在签名包含的时间戳之前，必须至少有1/4的**最低生命周期**内，没有进行过其他的创建新账号签名。
5. 两个执行签名的异属性账户前后进行签名，第二个签名的内容需包含第一个签名。
（*暂时并未强制两个异属性父级地址的签名顺序，但有可能在以后的自然法则定义中有所扩展，进一步提高创建账户的成本*）
6. 生命周期的长度跟父级账号相关，可定义为父级账户中生命周期较长的账号的1/2，但不低于最低生命周期。
（*关于生命周期的设定，有利于在整个系统诞生初始，通过对于时间证明的流速的控制，加速系统账号的扩展，同时在系统账号达到一定数量时，控制当前所有活跃账号的数量。*）
7. 子账户的生命周期，从执行签名的第二个父账户完成签名是，账户中包含的时间证明开始计算。
8. 两个执行签名的父级地址，不能为直接或间接的父子关系账户。
（*引入这个法则计算账户的公开地址过程中算法上的考虑，但也可以将其理解为，在创建新的账户过程中，我们必须要引入新的基因。*）

签名所产生的行为也属于一般消息（Message），其形式，传播方式及可信任程度均同于消息。

图例待补充……
 
## Message

消息（Message）在PDU中特指一个可包含信息的结构形式，一个消息可以包含一个或多个消息，而且嵌套的层级没有限制。整个分布式系统之中，账号之间的所有信息交互都是基于消息完成的。

消息结构，图例待补充……

#### Message Credit 

每个账户所生产的信息是否被接收，或者说其传播程度完全由接受者决定，而并非存在一个特定的中心化三方平台来保证消息的传播，也不存在共识机制来保证整个P2P的网络都认可接收某个信息，或者拒绝接收某个信息。对于是否接收某个新的消息的判断依据，来源于自身主观认定的消息来源账号的信用度和那条消息本身的可信程度及实时性。而信息的生产者，为了达到自己所生产的信息有更大传播概率的目的（*根本目的*）也会在构建消息的时候，尽量去符合PDU系统规则，以提高消息可信度。

在PDU中，下列的消息按照可信程度由弱到强进行排列：
1. 纯内容消息
2. 内容 + 数字签名
3. 内容 + 工作量证明 + 数字签名
4. 内容 + 时间证明/多维度时间证明 + 工作量证明 + 数字签名
5. 以链的形式维护自身生产的所有消息，且消息有序。

在整个系统的消息传播过程中，我们推荐终端用户只接收可信度为上述4，5的两种情况(除创世区间内)。按照4的可信级别构造消息的场景更多的是用于信息的转发，而5.的可信级别则作为正常信息发布。而且，如果某个第三方服务以某个账户为身份提供服务，则推荐按照5的信用级别来构造所有消息。总而言之，一个消息越容易被证伪，则其在没有被证明做恶之前则越可信。

#### Time Proof Message

虽然我们推荐尽量选取以最高信用等级方式传播的消息为时间证明，但理论上任何一条消息都可以作为时间证明，包含在你的消息当中。通常，在创建消息时，消息将包含的消息列表中的第0个位置作为时间证明；如果不需要时间证明，则第0个位置为空；如果需要多个时间证明，则在位置0包含一个列表的结构，其中包含多条时间证明消息。

#### Cosigned Birth Message

创建账户的过程中，生成新账户基本信息的过程通常不会被构造成消息在PDU中传播，因为此时待建账户并不合法，其他的账户不会接受此类消息。cosign过程所需的两个账户，有先后顺序，系统只要求后签名的父级地址将待建账户信息和叠加了两次签名的内容构造成消息（Message），在网络进行广播，第一个签名的地址不必须广播签名消息。但因为两次签名都必须带有时间证明，所以即便某账户在创建账户的过程中为先签名的账户且并没有发出过消息，如果被发现其两次创建账户的签名时间，小于1/4个最低生命周期，依然会被作为证据消息（Evidence Message）进行广播并处罚。

关于同一个公钥被多个私钥分别签名的情况，系统中也是允许的，相当于创建了多个同密码账户。

#### Evidence Message

证据通常为两条或条以上的矛盾信息，账户可以将违反自然法则的矛盾信息，包含到一条消息中进行转发。任何接收消息的账户都可按照自己的标准对于发出违反自然法则的账户处以处罚（不是转发证据消息的账户）。


## Network

#### Message Spread

通用的完整消息传播过程有如下几个步骤：
1. 生产消息内容，内容既可以由消息生产者自主产生，也可以是别的账号所生产并签名的内容。
2. 将时间证明和消息的内容进行组装，时间证明放在列表0的位置。
3. 为由2生成的消息体添加本账号前一条消息的hash，并添加签名，放入local账户的blockchain的头部。
4. 向临近的节点广播自己生产的消息，也可以只向目标的账户地址发送本条信息。（*临近的节点有账户的地址所定义，并非DAG关系图谱上所指的亲源关系。*）
5. 某账户地址收到消息后验证消息签名，如果签名的账户不在本地的DAG关系拓扑中，当前账户可以选择由其他账户地址请求信息完善自己的关系拓扑之后再处理本条消息，也可以直接放弃本条消息。（*转发其他账户的消息时，通常已经将此账户维护到自己的亲源关系拓扑当中，可响应其他账户关于相关亲源关系信息的请求。*）
6. 对消息进行验证，判断时间证明等具体消息中的信息，并进行后续的处理。
7. 可选择的将消息对临近的节点进行广播。
8. 如果由于这条消息，发现某个违反自然法则的证据，则对于相关作恶的账户进行惩罚，如拒绝接受这个账户在此之后的消息。或者更严重的，可以根据关系图谱处罚相关的其他账户的。之后可以将证据对临近的节点进行广播。

#### Account Create

在P2P的环境中，为账户的创建添加成本是维护整个账户体系的基础，我们通过账户之间的亲源关系及添加新账户时必须满足的自然法则来实现这一目标。账户的创建过程如以下步骤：
1. 新账户A创建秘钥对，并将A的公钥提供给第一个符合签名条件的账户B进行签名。此过程通常不通过次PDU网络的消息系统，因为新账户A此时对其他账户来说还不合法。
2. 账户B对A的公钥进行签名之后，可以任何方式，将签名之后的消息提供给另一个签名账户C，C必须和B为异属性账户。
3. 账户C对B签名之后的消息进行签名，并将此消息进行广播。
4. 收到广播的节点，会验证B，C的cosign是否合法，如果合法，则将A添加到自身维护的账户关系拓扑图当中。（同样，如果不合法，则收集作恶证据并广播。）

#### PDU Evolve Step

PDU账户系统的创建发展过程通常会经历一下的几个步骤：
1. 构建创世文件，其中包含两个公钥（Adam，Eve），被认为是此PDU的账户系统拓扑图（DAG）的顶端，这两个账户的二元属性必须相异。
2. 创建初始的N代账户，此过程中，Eve（也可以是Adam）发布时间证明事件，Adam和Eve及其N代之内的子账户，根据Eve发布的时间证明，在符合自然法则的条件下，构建一定数量的账户。
3. 由上述的账户中某个账户开始启动一个稳定的时间证明服务器，以提供P2P环境下的最初时间证明服务，这是PDU的第一次分裂时空。在此之后，P2P环境中的用户就能够更方便的参与创建账户。
4. 账户系统中账户总数的增速度由时间证明来控制，时间流速越慢，账户总数增速越慢。
5. 出现多个不同时间流速的证明，PDU中产生时空分裂，用户按照自己的意愿选择时空（可同时符合多个时空）来创建新的账户，并使用。
6. 多时空并存。

待补充图……

## Function Node

一个节点（Node）通常指一个信息转发节点，跟账户没有硬性绑定的关系。一个账户可以同时通过多个节点来发布信息，一个节点也可以同时转发多个账户的信息；节点可以只提供单个时空的信息，也可以提供多个时空的信息。简单的说，我们所说的节点是一个可以响应请求提供消息的三方服务，类似于互联网中DNS服务。

#### Time Proof Node

系统中会存在多个时间证明服务器节点，节点上可以保存多个不同流速的时间证明的完整信息。账户可以从服务器上获取自己所在时空的最新时间证明，加入自己的消息当中。也可以获取某个时空的历史时间证明，用以验证第三方信息的合法性。

#### Account Node

针对单一或者多个时空，维护最新的，最完整的账户信息，包含合法的账户信息，收集账户的作恶证据等，帮助用户在接收到一个未知来源的消息时，完善本地的账户亲源关系拓扑图。

#### Tracker Node

维护节点的当前状态，是否在线，监听端口等信息，使得用户可以在P2P的环境下直接跟对方进行交互。

#### Message Node

收集和维护消息，每个消息节点都根据不同的主观意愿（算法）来决定自己所转发（广播）的消息内容。消息节点相当于当今互联网上的众多网站，区别是其中消息（内容）的归属权为消息生产者。



## Conclusion

本文中我们提出了一种去中心化的P2P社交网络的构想，通过引入时间证明，并以此为基础，给予P2P的网络上一切用户行为提供成本计算的依据。用户可以按照自己的意愿，选择自己所存在的网络时空（多个）。通过次方式，我们将用户的身份及用户的社交关系归还于用户本身，而非依存与某个特定的社交网络服务。

不同于以比特币为代表的去中心化数字加密货币，PDU并非用共识强制全网接收并维护唯一的一致数据，而是让用户按照自己的意愿去选择接收整个网络中对自己有意义的那部分信息。同时允许用户传播并收集作恶证据的方式，来进一步增加作恶的成本。同时DAG结构的关系维护，也可以方便的进行关联惩罚，进一步提高此成本。建议本地节点维护的数据除自身产生的消息之外，只有自相关的部分账号拓扑关系，相比于比特币节点需维护全量交易信息的方式，极大地提高了信息存储的效率。
 
一个（*已知的*）系统很难同时满足去中心化，效率及整体一致性，因为货币系统本身的特点，比特币选择了去中心化和整体信息的一致性，而根据信息传播的特点，PDU选择了去中心化和效率。我们认为在信息的传播过程中，单个节点无需实时获知全网的所有完整信息，也能够容忍由某账户恶意行为所产生的错误信息。
 
人类社会发展至今，就没有产生过任何一次全人类的共识。个体总在不停的选择自己所接受的价值观，且同时努力让自己的行为去更加符合那种价值观。不同的体制，朝代更替，信仰兴衰，就如同PDU世界中的一个个不同的时间证明……


