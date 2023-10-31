// Copyright © 2023 Rak Laptudirm <rak@laptudirm.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

#ifndef CHESS_ZOBRIST
#define CHESS_ZOBRIST

#include <cstdint>

#include "piece.hpp"
#include "square.hpp"

namespace Chess {
    class Hash {
        private:
            uint64_t internal = 0;

        public:
            constexpr inline Hash() = default;
            constexpr inline explicit Hash(uint64_t hash) : internal(hash) {}

            constexpr inline explicit operator uint64_t() const {
                return internal;
            }

            constexpr inline bool operator ==(const Hash&) const = default;

            constexpr inline Hash operator +(Hash rhs) const {
                return Hash(internal ^ static_cast<uint64_t>(rhs));
            }

            constexpr inline Hash operator +=(Hash rhs) {
                internal ^= static_cast<uint64_t>(rhs);
                return *this;
            }

            constexpr inline Hash operator -(Hash rhs) const {
                return Hash(internal ^ static_cast<uint64_t>(rhs));
            }

            constexpr inline Hash operator -=(Hash rhs) {
                internal ^= static_cast<uint64_t>(rhs);
                return *this;
            }
    };

    namespace Keys {
        constexpr Hash None = Hash(0);

        namespace {
            constexpr std::array<std::array<uint64_t, Square::N>, ColoredPiece::N> psqt = {{
                {0x083610fb1cd7c6a5, 0xa37f944be9dfc323, 0xf6abbe2515a93cbb, 0x014d5ce796d3ea21, 0x46762749c86b2be7, 0xaf8f7e5e5ed8dab6, 0x650f5e0808e360fa, 0x92392e42419e33d7, 0x3f00957bf619fabd, 0x277059f962b2ad51, 0xd5e6b582d55f02f8, 0x6a8fc1e493122621, 0xb93875281e1a9e10, 0xfdccfe46fd5c65b6, 0x8fe7670648261096, 0xfaf02033d4a8e4be, 0x4cdbf1c399a0d591, 0x15ab0047084d6a72, 0x04c803b639b31ccf, 0xafc8b6cdc9cd9178, 0x9f6489ce28d8e4df, 0x6e0f22474ea92533, 0xc67d7cfe40573fbc, 0xc6e2de374960b2d3, 0x3dd9ff4b4cb20377, 0x2732a77574a34c97, 0x90109f006eb02f00, 0xd1d6984031b00ea1, 0x2222761e1ff24f3c, 0x3046e312f5926dd8, 0x2ee49120253af727, 0x868f3eb27661d798, 0xb5c64ce3d8887ca5, 0xe7eb41a397897ef8, 0x8be01949fc53c6e3, 0xc431f31919856a9b, 0x427fea13e941741b, 0x545ac69f3d1c6634, 0x5330e8f007f7a79c, 0xe1017ea38e3edacc, 0x3fd71ac257d29c3a, 0x211161dd93d52f71, 0x4b828af57d3a4472, 0xb757239537eb85e1, 0x70594501903e1f99, 0xb29c35ab5d55ca77, 0xfee1f0e1793f9ae3, 0x1493c090bdf0e21d, 0xff558a38b78e694e, 0xb2f1501e42d8c37f, 0x52e51685a29c6033, 0xdf11a0bcc1c921d3, 0xa4517cced14456a7, 0xe8e7e7b5f94817a8, 0xe5e60a7e4c3153a6, 0x699fc03bfc3ad0b3, 0x3c07bb3c37d3d153, 0x6251bd8731c30cb2, 0xc3dea9c62c4edca8, 0x607c06832e583a9e, 0xa2574452c4b0dd15, 0xdd1b4c11b5a1ad7d, 0x04a2634682c1aaad, 0x8c165c27b93899a1},
                {0x7adfd3d554658027, 0xfd774b1530cf1356, 0xfbebe15b01385c83, 0x062d679429588cb4, 0x6752115c2c5326e8, 0x51b42635f0cdc9aa, 0xae93c5295995b5f8, 0xd7b0bcd44364a6c6, 0x3b5ff8aaa4b255a9, 0x6c7f1261a536649a, 0xe8aa5791cc441371, 0xd86b5875c7dcb86d, 0x9a46cfd78ed9b762, 0xa0e117135d96df38, 0x9478ea3e9293fb5a, 0x03a733f03155429c, 0xd693ff9c09f873e8, 0x2a3d8dad465630ca, 0x0edafa049fd439b0, 0x090729732b690837, 0x5279c76801154a6a, 0x005d1b1daadc0167, 0xe8460df1498fcf95, 0xc1f9c15076df65f5, 0x0e99df998d80d424, 0x82c9e119ed321b0a, 0xa8dba34133a2004c, 0x3bb2efc57cd90111, 0xf0ec0e4129421d3c, 0xc0782c93ad3142c5, 0xdd61e5b15ff6b122, 0x455dd5d93aed39d5, 0x43e84734883942a1, 0xf3e1b7621ac2f5f5, 0x2179dcc18a2e0bc3, 0xe53a1c459f32878b, 0xeba0a229f4d45afb, 0x7a8cfe54e35fc5e7, 0x036543ee6e22fe10, 0x95e5fffd0af43e20, 0xbbcb0800930bfb77, 0x9217dc6bb35ca3e6, 0xf2cb1ab44210a347, 0xc51cbb72992489db, 0xbef5df21c347a8e1, 0x11ab10dbdfb93abe, 0x2bc604b273b84e04, 0xb115232b2e73a311, 0x163477644bd47fb5, 0x4b254d8161f32805, 0x63ef3c964052f0f8, 0x98dff249223f96ca, 0x6b07106fd6bceddc, 0x768ff02e843aad10, 0xb577f171389c94bb, 0x366fbe11e18cee44, 0x26968ac24a683664, 0x5cf0f35aa2aa6bbf, 0xbb13cca6b6051c0a, 0xa8f18e41930fd83f, 0x2dd3abe39d4af1e3, 0xe5ef7fe684965153, 0xcf8485194d6cb250, 0xe4665a4568064f04},
                {0x28dfad0a205b2e9c, 0x3465686005390915, 0x3b90f6e1f6c56840, 0xe4109f19e9fa7f95, 0x11d46f28d3dace84, 0xfe2bb5b257be494f, 0x7c2967e1b1ed0b95, 0xe43b4a381a3a37cf, 0x695059d5ffe6fbbf, 0xb2f9e81b811a7170, 0xcf46e879c65fe0ad, 0xb9f97cd8a4d78595, 0xc02a516db8ae144f, 0xad435686fb04e9ec, 0xf82bbb6f352a3960, 0xe6e42dc57d2df3e0, 0xd187aa3cdbedd5b0, 0xf4aec79145d15fae, 0xae9c3fca7088fe8f, 0xf873076c70c5e238, 0x8e94cbcfbe2f8eb5, 0xa69dbbe1e61f1481, 0x57c6ac4cd8547a67, 0xee976d8cb38ecb47, 0x82c4c4591e6a3619, 0x2c17d11bfbfd153f, 0xd023af78940fafde, 0x09cb7b8b3635c0f8, 0x9d339b95075e5f21, 0x618d55829c196453, 0x99872d72aa4b5bb1, 0x28411a439cfab02f, 0x0447c4980dd18c0b, 0x0a727dd8203971a7, 0x4d64017ea28444f8, 0x7933f58f03881b90, 0x0408e8373ef716fa, 0x7cccc649e930bbad, 0x90af3b4043e9899c, 0x4c3d73f5fb212cb9, 0xaeb57acbe523727b, 0xce31b1ba42dfa5ec, 0xbb49d484582c2b00, 0x605e3e628c10baf6, 0x375b37391ac9f3e3, 0xcd9c35bf28764550, 0xf7fa103085c18847, 0x7515338408400c09, 0x68db9f000c9ae26d, 0x7ee7c64e4a40bac1, 0x5e4bfb864335d91b, 0x54460f903f65383c, 0x97d82484d05f13ba, 0xc2e48b075cc5ee40, 0x740dffe55366710c, 0xf625ead458cb5363, 0x25edad6808412086, 0x3c5f9a8f6b509e77, 0x0f45f0963da28643, 0xf1e7394e16dbad3d, 0x67aaffa8538ae041, 0xb9c83a569c2b2064, 0x623d092e66653e08, 0xaadd09b034e21dfe},
                {0x351b3cb6fa0afb17, 0xf3fa5057957e9f1f, 0x3caf5f931167c3a4, 0x0049d1915fd8ec1f, 0x8415b4cdb479775d, 0xe8c4292086c4105c, 0xa8bce7aee1239b7d, 0xfe39b02a48d2a9e0, 0xc739fe5dcd4457d3, 0x1403db8fb3519890, 0xe8b28db23ff09313, 0xbb5d403967d07997, 0xac490676033eff75, 0x16a04fa30d1bf9d3, 0x997217e09587296c, 0xf3117e27351004e4, 0x5d7f1450e6c84a24, 0x2bcdae26c841d5b9, 0x664feffb28482b8c, 0x493ecf1831366263, 0xe59b7e560c61528a, 0xc845abe4a1cba795, 0x002648c6bf4c69a8, 0xd3700303c87b0929, 0xb12fb9bb17affa29, 0x126230fb4c36768a, 0x2ff7d2f543443003, 0x7f9ea0aa559d889c, 0x937c4397b0a311d9, 0x624e3386c8bd3630, 0xcc7b2837959caa4a, 0x7a9895b2c073f315, 0x29269f35e4ff07c1, 0xb1724d353a0949d0, 0x5854240d00156398, 0xaac30e66022f4cd2, 0x3d573340cdc49599, 0xb61a17cc1d88375e, 0x2dbbb30344a74700, 0xe5961efe2fa46058, 0xbda64ec9369c19b5, 0x31c2ac9cf0309bc9, 0xccd07315b51b25dc, 0x4b8da2176d7ddf91, 0x564a16a24ca73266, 0x69b573ecad4ff466, 0x1e33e2e504f2aac3, 0x13ec566100843602, 0xf85ff42af43ab8e3, 0xf1e5f9f5acaec2ff, 0xc0268b39c159fe69, 0x2fa2016c847c3298, 0x23245f3213a20bf0, 0xa194b3e61730337e, 0xbcd2d5538f951936, 0x8af394651b992396, 0x4d8b850410bc371e, 0xfc6d20ee872a1778, 0x4e3bc79cace5ce19, 0x419dd7b26ff5cf91, 0xb86542be5df66369, 0x759ff91e508a169a, 0x2699351e889f4ea2, 0x9271485845fae691},
                {0xc3e6cbc2d58d54e2, 0x0c9d65764e662a03, 0x35398cf17f55e546, 0x36298d8994ef782f, 0x74a1686641906112, 0x932e26c31e2a841c, 0x742e57797e804b64, 0x8cd96f04c93bcd46, 0x8eaa7a1fb167256e, 0xb2b979d48293ced2, 0x148afc7b1ad4a2e2, 0xd6011dba4f25674b, 0xde9b1153c122b489, 0x971f14a615bea388, 0x634b1f6b0b3afb58, 0xd4aabc1364bb0003, 0x7e9b907828fde17f, 0xfc46a281078eb9fb, 0xc16d1a9dd6133f13, 0x5629856b3076ce38, 0xf712384f29bc651d, 0x715c38e6c60edae5, 0x41e21c89f20dec3d, 0x7016e3fabc4678d5, 0x01e0e17095413176, 0xbe802cac9b27004a, 0xe494c0ee82c3c208, 0x36beeaaf24f54f9d, 0x5566d05a46fb6521, 0xf36e57a275276137, 0x0b86532e3399794b, 0x4f36092bbcd8cf44, 0xe8657cf6ea841919, 0xc042d797999a1028, 0x955ed6e192c63428, 0x07567e07cd7066ad, 0x1096cbac96dc14dd, 0xdf0e1ae46713d10e, 0x829db5d6ee0fb300, 0xc5c539dfefb9bd54, 0x2f0fa6f16182da44, 0x9c97fbba009e51b8, 0x1735053fa6caab1a, 0x1d904c80cc2a0dcf, 0xfe2053329db48023, 0x0d866ad29a19b204, 0x463cb247f64d3b66, 0x2b64d2b61f3fb47a, 0x0808900fd4708fff, 0x3469cfbdd1bf9ee7, 0xc5418c0abbe1a5d6, 0x4de827479c338e12, 0x543c0d8641fc84b4, 0x7b6c8fb0111ebd02, 0xd3a2bb2a34ce1d44, 0xfb15c47f676ab7d7, 0x9e1f46ce9296ba13, 0x70aed462117ba0a4, 0xbf0b1eb5c6478634, 0x627c1d570c1527f5, 0x6783c93750818a46, 0x51d88b5799738381, 0x39c3ea29e83c603a, 0x231482df2f8d560d},
                {0xeff5eeb2a2b20b32, 0x48bb703400db90c5, 0xadee028408e7e3e8, 0x659a2e1b59c31f32, 0xee8881a63b2d62b5, 0xbd6d5581989bdd88, 0x6d531bdd223994f9, 0x776495a7d3403463, 0x33c8a19c4c5cc49e, 0xc69cfcedfe47ca25, 0xe8071dfa94c0413f, 0xd91e6c71a4a8a576, 0xd484d7e096b2d4d7, 0x07bff7a4a384d89b, 0x8c45618188fa0eee, 0x030326012537c059, 0xa0c2212939bde392, 0xb1d1dee94ec0650a, 0xb1a7eef0f841580f, 0x8da02c798c8e77b4, 0xa6aa60c55d25910d, 0xa2869d0f3c7c8636, 0x0858fb0b1be4b947, 0x215c03e88f12ab8d, 0x2c345d1776316fe2, 0xe25dadf27182eb8d, 0x1dce4c56d00834cf, 0xa38b7f785b4551ef, 0x9db3fb522619706e, 0x3de4776d073c1249, 0xef3cb77613dbb07a, 0xe57165c9708e6e5b, 0xae96b0e1485d60fa, 0x7cee5fe03af00323, 0x640e188aa7b52e44, 0xd315dad8edd4e988, 0x52ad94329655d1e2, 0xdf206e5499f2fd9f, 0x676a97d8dd036dc3, 0xc5abf94469845903, 0xb0c617d45824f4c1, 0x12c3420396ac6cf8, 0x3d0017d165733446, 0xcb20cf04679762d0, 0x939f82a3dfb029d7, 0x415ced5a648dc4d2, 0xcc0da63afddae269, 0x147d1ca927afe895, 0x39178fa5df6427a1, 0x6ff05d98ce3e0973, 0x6c6122ba5673a0ea, 0x43b79aa160e2b9f2, 0x83cff8354424a170, 0xf3afe5a144fdb94f, 0xa33ff2d730d0962f, 0x8b8aa9b1aa280114, 0xb241aa1f7b293b26, 0x0497eb0e482c1777, 0x761516f375dc62ef, 0x9ac971b4bc1da3af, 0x8e14e1927ff59bb5, 0x189bf5a0bb82a62f, 0x73327c05cb3009a3, 0x9655c388016c3fe2},
                {0xa38152e5792c41dd, 0x262270c3737300b1, 0x33b1082ff0c8e331, 0x8eea7c34adae9a6d, 0x95230505c46b9a3d, 0xde8f0350047fb7a6, 0xf41592ec09662620, 0x5f7daa8e72708b86, 0x07c6fe7d5a169624, 0x5bf5ae615cd3bf25, 0x250eee0284fd0950, 0x3b673e349479cbee, 0x145f4ed31313bfc4, 0x69c026f532c3d433, 0xb946085d9a96daf2, 0x8cb2f1089fe5c7bd, 0x5e2c8d1ab19db4bf, 0x379b61b49d3525e0, 0xf344242925559c19, 0x1f558fc5ea7eb9be, 0xe2e8f392da038fe8, 0xb188b13b69086ca9, 0xd659336635ed6e74, 0x352a293989b52bdd, 0xf25988bb0b15c76e, 0xd032c19a0604d849, 0xf55dce120e5b70da, 0x0508c99da18984fd, 0x245ea813e90f9f7f, 0x96f24024ea008b2c, 0xcc115c56313a9d69, 0x74294f3b06a8833d, 0xea90ac815b457e75, 0x41649127eb1c4ce9, 0x20689236e3a8871e, 0xd678cfd8f1332076, 0x53d0414c27c5be8e, 0x49fb49539f3f4011, 0x5efb7f5936d930cc, 0xd06ce79c4ee00ca3, 0x517607ed03a758c9, 0x857f0d52e12edfa0, 0x620c0fbb2d6efc58, 0xc3780c4225407b19, 0xf62c4f10f9ecd54e, 0xfd9b6353aa8e64ca, 0xde268ff6dc85969c, 0x3c0bdb4f34b27e27, 0xa24a1ef85b4edaa9, 0xdb1f35914fc30fe9, 0x785a1b1a28468f79, 0x54cac7eb27f16f29, 0x5699b8193713e404, 0xf4f41920939d2f09, 0xbd3c0939d538f5bf, 0xee67fb624d3f279a, 0x0993bafa486dbfd0, 0x0bbfb4f7f6017912, 0x9eba8ece3a5e0aed, 0x0e93cfff50edec0a, 0x91844c5094791de6, 0xb240871946900373, 0x5a15f04e16e336f3, 0xae8506b7e0178da8},
                {0xcf1c140354d90d8d, 0xff011f11a27e1db5, 0x2f81119b6645bef5, 0xd3a5f1bcc336ef9a, 0xd09c41011c888ab4, 0xd6342e300e40c410, 0x577eb38e32439a91, 0xb16ffd8e6ede433f, 0x88201e51dbca9b91, 0x87c7b999dc878b73, 0xfbb96e76d739caf2, 0xffc91f5554e883f7, 0xfbdb1bb1163963e1, 0xb033e55a5bff12e9, 0x19bdbbe311bbfe5a, 0xb28c6c7c5f400188, 0xd8fecbcf3e92ee98, 0xf11abdf07f1033e4, 0x22a2fc6307fcdeb9, 0x9c180ffc0e3fb854, 0xedbca52dad4d07ed, 0x9e868776493703df, 0x1622a29ac26dc40e, 0x361f1333383764fe, 0xd6b1f3a9caa1ed2e, 0x23b335f0cb796d16, 0xc64a4d902a8f0661, 0x37fdfae72d1b30bc, 0x323ae9bd68fe607b, 0xae5e7beceb4953ff, 0x5b179e4261ab93af, 0x220eeb559046a5d2, 0x01b4229f83c1a79c, 0x39264dd39d1eea01, 0xbfdd7bfdb2a9e9ea, 0x3426f3b421450242, 0x2e77bc017c10cfa8, 0x99d60f361847d387, 0x42806cbdbbc55504, 0xe85708e048659f06, 0xbc132fd0e2e0976a, 0xfa686efea79c6da5, 0xfd058cb748ea808e, 0xee2d992c2f806e6f, 0xf9569c53380f7d24, 0x3943d426426ea766, 0x6ac6af3dd5df17f2, 0x6cde51169d69e52c, 0xd28b5d4c62d479ca, 0x4404dc78f30923eb, 0xa04c03f4a0f58b3a, 0x773c0f09934e0620, 0x5bcaa56f3bfe4271, 0xd950fbb6b80b7ce6, 0x73ab5233e3c02dbe, 0xc67fb2836190b3e3, 0xfc60852ab1bdeb2f, 0x8aee110872e49998, 0x555ed5746bbe8727, 0xdd6f1888daed759c, 0xcc5c915267ab26ba, 0x7de30f97853b00ac, 0x3b3cf0b03e3654d8, 0x348fec5cc59b0497},
                {0x3011c4d28635dbdf, 0x13b174f3eefdc297, 0x41c1aa861dc79560, 0x96fff72f157413d6, 0x546e8e8ec8773076, 0xd5b58b684d1a5399, 0x8bdb03e3e6d29838, 0x421c53655bbc1521, 0x1c920a8701f626cf, 0xe172bfb282e929b1, 0xae27d629badb1b6d, 0x4738ec83a85f112a, 0xb7566e63c52f73ff, 0x6fb5e187fbd0757e, 0xc52fc3ed8ff08176, 0xd03bb85163751086, 0x258aaa40c155846d, 0x5bb09b8ea743858a, 0x7d707997049f506a, 0x88e5c579e8b8ec8f, 0x7170a24e2c0c8a00, 0xdee1d4843e7d7907, 0x4c1e766b2ee31c35, 0xacdf4cca41fd08af, 0x7bc78d0083b84854, 0xd71eff4935d3c228, 0x2d01451ad4d06582, 0x523d9682a4d37017, 0x58e39191f3cb587a, 0x026515714520fc53, 0xeffaa5630885430d, 0xbadaba2091156ac1, 0x33277e8b0439291e, 0x7aea720c476f6645, 0xd605947274c6cd23, 0x34f4d8e26e91bb5e, 0x2fa33797aee09da6, 0x0b5b426be0430939, 0x3880f1f85a0f6ab4, 0xb882fc47309805fa, 0x21aceae54062f31f, 0x8bd6386fc481372e, 0x79e7b84b6f039893, 0x299820e9679f0906, 0xdadbb60cb96722d4, 0xb4a69d5a5125f3ad, 0x3c1a02477403c485, 0x97bf24886211b282, 0x8fb9f64dd9c7e655, 0x1d1e7319dce7412f, 0xcd3eacf88a4ce2c4, 0x9c251f9570f4a41e, 0x6440d17499eba25d, 0xd0b507d56ae36045, 0xb766d402e56f0d8d, 0x144b20dca1156997, 0x4fed16b58e4b6e2b, 0x4ff60ff14a592e41, 0x1b049bdea4d05426, 0x79d6502120c6c8e1, 0x8a810ff080a3e083, 0x7d26ed2c1eb6ebc5, 0x8d371c46110d0b72, 0xf53957ac0caab20c},
                {0x1c6a15e74c484818, 0x394ecc7315c776b3, 0x8b338c025467af83, 0x755df72e74e28c2c, 0x096102c2f4721596, 0xda324813d5f5165c, 0x13a72cf0f2f0c8c4, 0xfe8772410008712a, 0x3b640efbb53b4127, 0x69779f11fc633452, 0x75de90b625fda51d, 0x4b9c82ee1e1cb305, 0x6eace48f276be344, 0x32d00fceb789ec71, 0xf1faa8b8a4addd4a, 0x6b2dd36fbf2e5ec4, 0xad2bb7a46b82cab4, 0x49012620972ce6ce, 0x32dc03c3cf95b8b8, 0xa9f463724298da92, 0x9e80e8729b9e098e, 0x94a5f1293de1972c, 0x0577e33a55f297eb, 0x16f3b7b1b2c800d4, 0x934d62300037b090, 0x30ba5035eaa9f1d3, 0xcdca15d562592c40, 0xb0aae4af24edd99d, 0x7eb866dc206dfa52, 0x91602ec574b77474, 0xa98abd14dde57859, 0xaef082e17aae0e3d, 0x00c39cb0f82e24a1, 0x4ea8d7b26183d512, 0x49d058a520fcfcfc, 0x50a8f5a501b860ff, 0xac97a5b426ab824a, 0x9efc8ac042139f45, 0xf0d84b3d42b5cb99, 0xb1e8c0adab3d57d7, 0x1c7a0fba85a8aaeb, 0x87565f24bdc3ee7a, 0x77552ed09b8b4101, 0x95ee84237775535c, 0xf148623c65791a53, 0x306f04eadff39f55, 0xcfb27c101bfc3dae, 0x25b1bc975e125ba6, 0xbe2e97660e85f62b, 0x55350c3c99bb7a26, 0xa72aba5099663783, 0x5198c5e6a82368d3, 0xfe68bbdf927faa6e, 0x7338bf90c9ed7039, 0x2e5078e9d6b3b8e5, 0x40684cd6b9c6cac0, 0xf3979178e731c738, 0xd392f50ab651e966, 0x0c7916677a67f9aa, 0xbac5b81b53946b68, 0xf47d692e0a0ae20e, 0xaf98a3b93ac483fd, 0x36c3343929a28281, 0x01177bbc613bdfd7},
                {0x68085e26dbe3ad56, 0x9a9d46582a40120b, 0x8aa6abbd2cad7d96, 0x5527a24035773ed8, 0xc79805af15fa519c, 0xa9a03e8fb9f60885, 0x82f999d825db04e0, 0x49db5f367e106034, 0x83fbfc6a4aa8f161, 0xc1daedbaa5d01451, 0x7d938e607492dfe8, 0x622135de5b37f9c1, 0x6946d729ce3a1019, 0xb19a3dfdd10d34a8, 0xbff22fd4f4268351, 0xc329a8b2c951b7ff, 0x63da62e7e591dcec, 0xbf007b12ec4307ac, 0x792444890a0570c6, 0x72318d01e4ccf0a4, 0x50e0d2417bdb719b, 0x1565a2897030890e, 0xf9d5d18956242293, 0x64104ef221973e5a, 0x5dd2fda8c41eb447, 0x175ed04f5cba8520, 0x4b41274dc059c1de, 0x52c6a011722f7525, 0xdeb942504bc8e782, 0xb458d3594d6cae08, 0x1eac4cb3fa22358e, 0xb8b970f1500a1119, 0x3c74e78cc4a6420f, 0x978ef947dd452dce, 0x3e2e951e6b2f0efe, 0xa56f9e5d36f3a00b, 0xf77371e0e30687d4, 0xf530ae19bf5498e5, 0x772163240b406f47, 0x8bf14ec5102856f2, 0xd29afaf89fbc4012, 0x2f37b6297c95b3f0, 0xf99323223fa8d818, 0xbd33ffd00a14c9aa, 0xfc8af274e35822fe, 0x635a69eaa68adee7, 0x57d645d580e935f0, 0x3fc98238def97d41, 0x1ac557171e66091b, 0x28d6dd4d2a8e542c, 0xf47a8200e4b78fa8, 0xcb27461f07dcaeda, 0x0344565cd7c80558, 0xd6f32dd8e7a4c265, 0xc963e291da80d2ff, 0x441d93cafd5df3df, 0x6f0df8634290aa45, 0x0556b564010e6b21, 0x3d3e34e8eff6e213, 0xdf37a92c959fc1b8, 0x6c7c380625981e73, 0x9fe365590db2e003, 0x9391b03d2f536994, 0x6188e8d1db75331d},
                {0xabab879cd5585f2f, 0xfdb8a69bc4052dd5, 0xa097af8b98ae5653, 0xa7262be7fa75d97b, 0xda8f8ae4c5526fba, 0xac8d445dc93990b3, 0x311e44664ea37966, 0x72358b3b76d6e28b, 0xfd84b139d74da2ad, 0xfbad215ccd898848, 0x8c7a00a136a05ffd, 0x7709e685c945ee73, 0xeb32efd0627aecc1, 0x3e6f41983f953cd8, 0x46ebf3bd647cc189, 0x21e91003e0e722b7, 0x5ff78aee36f5e7df, 0x7f0b0b2514024f0f, 0x31a7b80fad47192f, 0xd48ca8c3be089ea4, 0x6220c3ea0477a100, 0xcda3d82077f85837, 0x29a7477b3274955b, 0xb46b8fa6c96a547c, 0xc76e82f848d82a29, 0x9912a9640c62023d, 0xc59e8a1a77cabde7, 0x82ac3fd8bb87ecff, 0x5c7fb3bfff378cbb, 0xb0a9a087ea30e56f, 0x01c4f4855092269f, 0x53e0dc61631cfd20, 0xb482604ea6d2a918, 0xc0be737023dcdef6, 0xbbdb426b8e95919e, 0xe4e54404356b9992, 0x1d8fd20388787282, 0x4a85dc29bf8e1109, 0x450eb0cb187bcafb, 0xf51e953f2053516a, 0x8d7a82dfecd6f2f0, 0x82ee9c1328eaf825, 0x80b8a490de34e58c, 0xc199c2cf6fa3c4a0, 0x404f57fd165644eb, 0xf335001fc9324ab4, 0xb1109adca3c18129, 0x2b65dc52c43442c5, 0x36f814c72a173952, 0xce5c402e9cf3bc46, 0x043c3cba93773393, 0x397305568e833188, 0x03c8b53be7ebb8f4, 0xd8c9ea4dbbe0caba, 0xe4c12637188a7f2f, 0xb3c39c29782b86c8, 0x9430009ef3092669, 0xfa7d3f1cc2dae40e, 0x6ead2df26cbef22b, 0x92060073bd794085, 0xaef2c95bd9ad5886, 0xc13f07c270b5cace, 0x5b21dd821267ea79, 0x2fe9a4d5aa8d43f6},
            }};

            constexpr std::array<uint64_t, File::N> ep = {
                0x14c6099d731723b7,
                0x1cec25e490795dfb,
                0xa2c8015acdd7305f,
                0xc65d7c2700f3aade,
                0xe0fe6bcd9c147fb1,
                0x593b8aea38433907,
                0x2fe646b777886e9f,
                0xc045c1dde772ac79,
            };

            constexpr uint64_t whiteH = 0x4d28598573750b10;
            constexpr uint64_t whiteA = 0xdfe34de8892603ad;
            constexpr uint64_t blackH = 0x177ab8314c2b200e;
            constexpr uint64_t blackA = 0xc07e0a697776ea93;

            constexpr std::array<uint64_t, 16> castling = [](){
                std::array<uint64_t, 16> keys = {};

                for (uint8_t rawRights = 0; rawRights < 16; rawRights++) {
                    keys[rawRights] = 0UL;

                    auto rights = Castling::Rights(rawRights);
                    if (rights.Has(Castling::WhiteH)) keys[rawRights] ^= whiteH;
                    if (rights.Has(Castling::WhiteA)) keys[rawRights] ^= whiteA;
                    if (rights.Has(Castling::BlackH)) keys[rawRights] ^= blackH;
                    if (rights.Has(Castling::BlackA)) keys[rawRights] ^= blackA;
                }

                return keys;
            }();
        }

        constexpr Hash SideToMove = Hash(0x5ec3a196160b9a06);

        constexpr inline Hash PieceOnSquare(ColoredPiece piece, Square square) {
            return Hash(psqt[static_cast<uint8_t>(piece)][static_cast<uint8_t>(square)]);
        }

        constexpr inline Hash EnPassantTarget(Square target) {
            return Hash(ep[static_cast<uint8_t>(target.File())]);
        }

        constexpr inline Hash CastlingRights(Castling::Rights rights) {
            return Hash(castling[static_cast<uint8_t>(rights)]);
        }
    }
}

#endif
