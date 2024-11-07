#include <iostream>
#include <fstream>
#include <sstream>
#include <llvm/IR/LLVMContext.h>
#include <llvm/IR/IRBuilder.h>
#include <llvm/IR/Module.h>
#include <llvm/IR/Value.h>
#include <llvm/Support/raw_ostream.h>
#include <map>
#include <string>
#include <memory>

using namespace llvm;

static std::unique_ptr<LLVMContext> TheContext;
static std::unique_ptr<IRBuilder<>> Builder;
static std::unique_ptr<Module> TheModule;
static std::map<std::string, Value *> NamedValues;

void compile(std::string program)
{
    TheContext = std::make_unique<LLVMContext>();
    TheModule = std::make_unique<Module>("my cool bfcompiler", *TheContext);
    Builder = std::make_unique<IRBuilder<>>(*TheContext);

    ArrayType *TapeType = ArrayType::get(Type::getInt8Ty(*TheContext), 1024 * 1024 * 4);
    GlobalVariable *Tape = new GlobalVariable(*TheModule, TapeType, false, GlobalValue::PrivateLinkage,
                                              Constant::getNullValue(TapeType), "tape");

    GlobalVariable *tapeIndex = new GlobalVariable(*TheModule, Type::getInt32Ty(*TheContext), false, GlobalValue::PrivateLinkage,
                                                   Constant::getNullValue(Type::getInt32Ty(*TheContext)), "tapePointer");

    FunctionType *putcharType = FunctionType::get(Type::getInt32Ty(*TheContext), {Type::getInt32Ty(*TheContext)}, false);
    Function *putcharFunc = Function::Create(putcharType, Function::ExternalLinkage, "my_putchar", TheModule.get());

    FunctionType *getcharType = FunctionType::get(Type::getInt32Ty(*TheContext), {}, false);
    Function *getcharFunc = Function::Create(getcharType, Function::ExternalLinkage, "my_getchar", TheModule.get());

    FunctionType *mainFuncType = FunctionType::get(Type::getInt32Ty(*TheContext), false);
    Function *mainFunction = Function::Create(mainFuncType, Function::ExternalLinkage, "main", TheModule.get());
    BasicBlock *entryBlock = BasicBlock::Create(*TheContext, "entry", mainFunction);
    Builder->SetInsertPoint(entryBlock);

    for (char ch : program)
    {
        switch (ch)
        {
        case '>':
        {
            Value *currentIndex = Builder->CreateLoad(Type::getInt8Ty(*TheContext), tapeIndex, "currentIndex");
            Value *newVal = Builder->CreateAdd(currentIndex, ConstantInt::get(Type::getInt8Ty(*TheContext), 1), "newVal");
            Builder->CreateStore(newVal, tapeIndex);
            break;
        }
        case '<':
        {
            Value *currentIndex = Builder->CreateLoad(Type::getInt8Ty(*TheContext), tapeIndex, "currentIndex");
            Value *newVal = Builder->CreateSub(currentIndex, ConstantInt::get(Type::getInt8Ty(*TheContext), 1), "newVal");
            Builder->CreateStore(newVal, tapeIndex);
            break;
        }
        case '+':
        {
            Value *currentIndex = Builder->CreateLoad(Type::getInt8Ty(*TheContext), tapeIndex, "currentIndex");
            Value *ptr = Builder->CreateGEP(TapeType, Tape, {ConstantInt::get(Type::getInt8Ty(*TheContext), 0), currentIndex}, "tapePtr"); // ptr = TAPE + IDX
            Value *currentVal = Builder->CreateLoad(Type::getInt8Ty(*TheContext), ptr, "currentVal");
            Value *newVal = Builder->CreateAdd(currentVal, ConstantInt::get(Type::getInt8Ty(*TheContext), 1), "newVal");
            Builder->CreateStore(newVal, ptr);
            break;
        }
        case '-':
        {
            Value *currentIndex = Builder->CreateLoad(Type::getInt8Ty(*TheContext), tapeIndex, "currentIndex");
            Value *ptr = Builder->CreateGEP(TapeType, Tape, {ConstantInt::get(Type::getInt8Ty(*TheContext), 0), currentIndex}, "tapePtr");
            Value *currentVal = Builder->CreateLoad(Type::getInt8Ty(*TheContext), ptr, "currentVal");
            Value *newVal = Builder->CreateSub(currentVal, ConstantInt::get(Type::getInt8Ty(*TheContext), 1), "newVal");
            Builder->CreateStore(newVal, ptr);
            break;
        }
        break;
        case '.':
        {
            Value *currentIndex = Builder->CreateLoad(Type::getInt8Ty(*TheContext), tapeIndex, "currentIndex");
            Value *ptr = Builder->CreateGEP(TapeType, Tape, {ConstantInt::get(Type::getInt8Ty(*TheContext), 0), currentIndex}, "tapePtr");
            Value *currentVal = Builder->CreateLoad(Type::getInt8Ty(*TheContext), ptr, "currentVal");
            Builder->CreateCall(putcharFunc, currentVal);
            break;
        }
        case ',':
        {
            Value *inputVal = Builder->CreateCall(getcharFunc);
            Value *currentIndex = Builder->CreateLoad(Type::getInt8Ty(*TheContext), tapeIndex, "currentIndex");
            Value *ptr = Builder->CreateGEP(TapeType, Tape, {ConstantInt::get(Type::getInt8Ty(*TheContext), 0), currentIndex}, "tapePtr"); // ptr = TAPE + IDX
            Builder->CreateStore(inputVal, ptr);
            break;
        }
        case '[':
            // asm_b.Add_instr("cmp     BYTE [%s], 0", TAPE_PTR)
            // asm_b.Add_instr("je      right_%s", strconv.Itoa(bracket_pairs[i]))
            // asm_b.Add_label("left_%s", strconv.Itoa(i))
            break;
        case ']':
            // asm_b.Add_instr("cmp     BYTE [%s], 0", TAPE_PTR)
            // asm_b.Add_instr("jne     left_%s", strconv.Itoa(bracket_pairs[i]))
            // asm_b.Add_label("right_%s", strconv.Itoa(i))
            break;
        default:
            break; // do nothing
        }
    }

    Builder->CreateRet(ConstantInt::get(Type::getInt32Ty(*TheContext), 0));
}

int main(int argc, char *argv[])
{
    if (argc <= 1)
    {
        std::cout << "Usage: ./bf <bf-program.b>";
        return 0;
    }

    std::ifstream file(argv[1]);
    if (!file)
    { // Check if the file was successfully opened
        std::cerr << "Could not open the file" << argv[1] << std::endl;
        return 1;
    }

    std::stringstream buffer;
    buffer << file.rdbuf();

    std::string program = buffer.str();

    compile(program);

    // std::cout << program << std::endl;
    TheModule->print(llvm::outs(), nullptr);

    return 0;
}