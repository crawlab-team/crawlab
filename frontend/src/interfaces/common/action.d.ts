import {ACTION_CLONE, ACTION_DELETE, ACTION_EDIT, ACTION_RUN, ACTION_VIEW,} from '@/constants/action';

declare global {
  type ActionName =
    ACTION_VIEW |
    ACTION_EDIT |
    ACTION_CLONE |
    ACTION_RUN |
    ACTION_DELETE;
}
