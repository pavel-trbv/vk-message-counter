import { FC } from 'react';
import { Heading, Link, Text, chakra } from '@chakra-ui/react';

const Instruction: FC = () => {
  const tokenUrl =
    'https://oauth.vk.com/authorize?client_id=6121396&scope=69632&redirect_uri=' +
    'https://oauth.vk.com/blank.html&display=page&response_type=token&revoke=1';

  const TokenLink = (
    <Link href={tokenUrl} color="blue.500" target="_blank">
      this link
    </Link>
  );

  return (
    <>
      <Heading mb="4">How to use</Heading>
      <Text lineHeight="6" mb="4">
        To get message statistics, you need 2 things: api token and chat id.
        <br />
        <br />
        You can get the token by logging into {TokenLink} and copying part of the address bar from `access_token=` to `&expires_in`.
        <br />
        <br />
        To get the chat ID, you need to open the chat on the desktop version of vk.com and copy the numbers after `im?sel=c` from the address bar.
        <br />
        Example: https://vk.com/im?sel=c
        <chakra.span background="twitter.50" display="inline">
          123
        </chakra.span>
        .
      </Text>
    </>
  );
};

export default Instruction;
